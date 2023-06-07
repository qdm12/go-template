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

func (db *Database) Start() (runError <-chan error, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	runError, err = db.memory.Start()
	if err != nil {
		return nil, fmt.Errorf("starting memory database: %w", err)
	}

	err = db.initDatabaseFile()
	if err != nil {
		_ = db.memory.Stop()
		return nil, fmt.Errorf("initializing database file: %w", err)
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

	return db.writeFileNoLock()
}

func (db *Database) writeFileNoLock() error {
	const perms fs.FileMode = 0600
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

// initDatabaseFile either creates the database file or loads
// existing data from it into the memory database.
func (db *Database) initDatabaseFile() (err error) {
	exists, err := fileExists(db.filepath)
	if err != nil {
		return fmt.Errorf("%w: %w", dataerrors.ErrReadFile, err)
	}

	if !exists {
		return db.writeFileNoLock()
	}

	stat, err := os.Stat(db.filepath)
	if err != nil {
		return fmt.Errorf("%w: %w", dataerrors.ErrReadFile, err)
	}

	if stat.Size() == 0 {
		return db.writeFileNoLock()
	}

	err = db.readFile()
	if err != nil {
		return fmt.Errorf("%w: %w", dataerrors.ErrReadFile, err)
	}
	return nil
}

// readFile is only used when initializing the database.
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
