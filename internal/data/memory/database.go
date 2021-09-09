// Package memory implements a data store in memory only.
package memory

import (
	"sync"

	"github.com/qdm12/go-template/internal/models"
)

// Database is the in memory implementation of the database store.
type Database struct {
	sync.RWMutex
	data models.Data
}

// NewDatabase creates an empty memory based database.
func NewDatabase() (*Database, error) {
	return &Database{}, nil
}

func (db *Database) Close() error {
	db.Lock()
	defer db.Unlock() // wait for ongoing operation to finish
	return nil
}

func (db *Database) GetData() models.Data {
	db.Lock()
	defer db.Unlock()
	return db.data
}

func (db *Database) SetData(data models.Data) {
	db.Lock()
	defer db.Unlock()
	db.data = data
}
