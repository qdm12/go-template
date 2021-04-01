// Package psql implements a data store using a client to a
// Postgres database.
package psql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/qdm12/go-template/internal/config"
	"github.com/qdm12/go-template/internal/data/errors"
	"github.com/qdm12/golibs/crypto/random"
	"github.com/qdm12/golibs/logging"
)

// Database is the Postgres implementation of the database store.
type Database struct {
	sql    *sql.DB
	logger logging.Logger
	random random.Random
}

// NewDatabase creates a database connection pool in DB and pings the database.
func NewDatabase(config config.Postgres, logger logging.Logger) (*Database, error) {
	connStr := "postgres://" + config.User + ":" + config.Password +
		"@" + config.Address + "/" + config.Address + "?sslmode=disable&connect_timeout=1"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrCreation, err)
	}
	fails := 0
	const maxFails = 3
	const sleepDuration = 200 * time.Millisecond
	for {
		err = db.Ping()
		if err == nil {
			break
		}
		fails++
		if fails == maxFails {
			return nil, fmt.Errorf("%w: %s", errors.ErrCreation, err)
		}
		time.Sleep(sleepDuration)
	}
	return &Database{db, logger, random.NewRandom()}, nil
}

// Close closes the database and prevents new queries from starting.
func (db *Database) Close() error {
	if err := db.sql.Close(); err != nil {
		return fmt.Errorf("%w: %s", errors.ErrClose, err)
	}
	return nil
}
