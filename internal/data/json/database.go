// Package json implements a data store using a single JSON file
// and the memory package.
package json

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/qdm12/go-template/internal/data/errors"
	"github.com/qdm12/go-template/internal/data/memory"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/golibs/files"
)

// Database is the JSON file implementation of the database store.
type Database struct {
	sync.Mutex
	memory      *memory.Database
	filepath    string
	fileManager files.FileManager
}

// NewDatabase creates a JSON file at the filepath provided if needed,
// and reads existing data in memory.
func NewDatabase(memory *memory.Database, filepath string) (*Database, error) {
	db := Database{
		memory:      memory,
		filepath:    filepath,
		fileManager: files.NewFileManager()}
	exists, err := db.fileManager.FileExists(filepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrCreation, err)
	} else if !exists {
		if err := db.fileManager.Touch(filepath); err != nil {
			return nil, fmt.Errorf("%w: %s", errors.ErrCreation, err)
		}
	}
	rawData, err := db.fileManager.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrCreation, err)
	} else if len(rawData) == 0 {
		if err := db.writeFile(); err != nil {
			return nil, fmt.Errorf("%w: %s", errors.ErrCreation, err)
		}
	} else if err := db.readFile(); err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrCreation, err)
	}
	return &db, nil
}

func (db *Database) Close() error {
	if err := db.memory.Close(); err != nil {
		return err
	}
	db.Lock()
	defer db.Unlock() // wait for ongoing operation to finish
	return nil
}

func (db *Database) writeFile() error {
	db.Lock()
	defer db.Unlock()
	b, err := json.Marshal(db.memory.GetData())
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrWriteFile, err)
	}
	if err := db.fileManager.WriteToFile(db.filepath, b); err != nil {
		return fmt.Errorf("%w: %s", errors.ErrWriteFile, err)
	}
	return nil
}

// readFile is only used when creating the database.
func (db *Database) readFile() error {
	b, err := db.fileManager.ReadFile(db.filepath)
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrReadFile, err)
	}
	var data models.Data
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("%w: %s", errors.ErrReadFile, err)
	}
	db.memory.SetData(data)
	return nil
}
