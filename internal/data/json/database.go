// Package json implements a data store using a single JSON file
// and the memory package.
package json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

// NewDatabase creates a JSON Database object with the memory
// database and filepath given. Its `Start` method will either
// initialize the JSON database file or load existing data from
// an existing JSON file into the memory database.
func NewDatabase(memory *memory.Database, filepath string) *Database {
	return &Database{
		memory:   memory,
		filepath: filepath,
	}
}

func (db *Database) String() string {
	return "JSON file database"
}

func (db *Database) Start(ctx context.Context) (runError <-chan error, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	runError, err = db.memory.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting memory database: %w", err)
	}

	err = db.loadDatabaseFile()
	if err != nil {
		_ = db.memory.Stop()
		return nil, fmt.Errorf("loading database file: %w", err)
	}

	return runError, nil
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

	const perms fs.FileMode = 0o600
	file, err := os.OpenFile(db.filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perms)
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

// loadDatabaseFile loads the data from the database file
// if the file exists and is not empty. If the file does not
// exist, its path parent directory is created.
func (db *Database) loadDatabaseFile() (err error) {
	file, err := os.Open(db.filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			const perm fs.FileMode = 0o700
			return os.MkdirAll(filepath.Dir(db.filepath), perm)
		}
		return fmt.Errorf("%w: %w", dataerrors.ErrReadFile, err)
	}

	stat, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("%w: %w", dataerrors.ErrReadFile, err)
	} else if stat.Size() == 0 { // empty file
		_ = file.Close()
		return nil
	}

	decoder := json.NewDecoder(file)
	var data models.Data
	err = decoder.Decode(&data)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("%w: %w", dataerrors.ErrDecoding, err)
	}
	db.memory.SetData(data)

	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing database file: %w", err)
	}

	return nil
}
