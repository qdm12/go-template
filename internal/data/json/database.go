// Package json implements a data store using a single JSON file
// and the memory package.
package json

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"sync"

	dataerrors "github.com/qdm12/go-template/internal/data/errors"
	"github.com/qdm12/go-template/internal/data/memory"
	"github.com/qdm12/go-template/internal/models"
)

// Database is the JSON file implementation of the database store.
type Database struct {
	mutex    sync.Mutex
	memory   *memory.Database
	filepath string
}

// NewDatabase creates a JSON file at the filepath provided if needed,
// and reads existing data in memory.
func NewDatabase(memory *memory.Database, filepath string) (*Database, error) {
	db := Database{
		memory:   memory,
		filepath: filepath,
	}
	exists, err := fileExists(filepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dataerrors.ErrReadFile, err)
	} else if !exists {
		const perms fs.FileMode = 0600
		err = os.WriteFile(filepath, nil, perms)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", dataerrors.ErrWriteFile, err)
		}
	}
	rawData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dataerrors.ErrReadFile, err)
	} else if len(rawData) == 0 {
		if err := db.writeFile(); err != nil {
			return nil, fmt.Errorf("%w: %w", dataerrors.ErrWriteFile, err)
		}
	} else if err := db.readFile(); err != nil {
		return nil, fmt.Errorf("%w: %w", dataerrors.ErrReadFile, err)
	}
	return &db, nil
}

func (db *Database) String() string {
	return "JSON file database"
}

func (db *Database) Start() (runError <-chan error, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.memory.Start()
}

func (db *Database) Stop() (err error) {
	err = db.memory.Stop()
	if err != nil {
		return fmt.Errorf("stopping memory database: %w", err)
	}
	db.mutex.Lock()
	defer db.mutex.Unlock() // wait for ongoing operation to finish
	return nil
}

func (db *Database) writeFile() error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	const perms fs.FileMode = 0600
	file, err := os.OpenFile(db.filepath, os.O_WRONLY|os.O_TRUNC, perms)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(db.memory.GetData())
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("encoding data to file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing file: %w", err)
	}
	return nil
}

// readFile is only used when creating the database.
func (db *Database) readFile() error {
	file, err := os.Open(db.filepath)
	if err != nil {
		return fmt.Errorf("%w: %w", dataerrors.ErrReadFile, err)
	}

	decoder := json.NewDecoder(file)
	var data models.Data
	err = decoder.Decode(&data)
	if err != nil {
		return fmt.Errorf("%w: %w", dataerrors.ErrDecoding, err)
	}
	db.memory.SetData(data)
	return nil
}
