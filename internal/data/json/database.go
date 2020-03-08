package json

import (
	libjson "encoding/json"
	"fmt"
	"sync"

	"github.com/qdm12/golibs/files"

	"github.com/qdm12/REPONAME_GITHUB/internal/data/memory"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
)

// Database is the JSON file implementation of the database store
type Database struct {
	sync.Mutex
	memory      *memory.Database
	filepath    string
	fileManager files.FileManager
}

// NewDatabase creates a JSON file at the filepath provided if needed,
// and reads existing data in memory
func NewDatabase(memory *memory.Database, filepath string) (*Database, error) {
	const errorWrapper = "cannot create JSON database"
	db := Database{
		memory:      memory,
		filepath:    filepath,
		fileManager: files.NewFileManager()}
	exists, err := db.fileManager.FileExists(filepath)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", errorWrapper, err)
	} else if !exists {
		if err := db.fileManager.Touch(filepath); err != nil {
			return nil, fmt.Errorf("%s: %s", errorWrapper, err)
		}
	}
	rawData, err := db.fileManager.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", errorWrapper, err)
	} else if len(rawData) == 0 {
		if err := db.writeFile(); err != nil {
			return nil, fmt.Errorf("%s: %s", errorWrapper, err)
		}
	} else if err := db.readFile(); err != nil {
		return nil, fmt.Errorf("%s: %s", errorWrapper, err)
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
	b, err := libjson.Marshal(db.memory.GetData())
	if err != nil {
		return fmt.Errorf("cannot write data to JSON file: %s", err)
	}
	if err := db.fileManager.WriteToFile(db.filepath, b); err != nil {
		return fmt.Errorf("cannot write data to JSON file: %s", err)
	}
	return nil
}

// readFile only used when creating database
func (db *Database) readFile() error {
	b, err := db.fileManager.ReadFile(db.filepath)
	if err != nil {
		return fmt.Errorf("cannot read JSON file: %s", err)
	}
	var data models.Data
	if err := libjson.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("cannot read JSON file: %s", err)
	}
	db.memory.SetData(data)
	return nil
}
