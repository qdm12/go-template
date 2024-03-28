// Package memory implements a data store in memory only.
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/goservices"
)

// Database is the in memory implementation of the database store.
type Database struct {
	sync.RWMutex
	data    models.Data
	running bool
}

// NewDatabase creates an empty memory based database.
func NewDatabase() (*Database, error) {
	return &Database{}, nil
}

func (db *Database) String() string {
	return "memory database"
}

func (db *Database) Start(_ context.Context) (runError <-chan error, err error) {
	db.Lock()
	defer db.Unlock()
	if db.running {
		return nil, fmt.Errorf("%w", goservices.ErrAlreadyStarted)
	}
	db.running = true
	return nil, nil //nolint:nilnil
}

func (db *Database) Stop() (err error) {
	db.Lock()
	defer db.Unlock() // wait for ongoing operation to finish
	if !db.running {
		return fmt.Errorf("%w", goservices.ErrAlreadyStopped)
	}
	db.running = false
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
