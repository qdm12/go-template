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

// PeriodicHealthcheck pings the database periodically and runs onFailure if an error occurs.
func (db *Database) PeriodicHealthcheck(period time.Duration, onFailure func(err error)) {
	db.RunTaskPeriodically(period, func(db *Database) {
		if err := db.sql.Ping(); err != nil {
			onFailure(err)
		}
	})
}

// RunInitialTasks runs tasks asynchronously and waits for all
// of them to complete.
func (db *Database) RunInitialTasks(tasks ...func(db *Database) error) (errors []error) {
	chError := make(chan error)
	for _, task := range tasks {
		task := task
		go func() {
			chError <- task(db)
		}()
	}
	for range tasks {
		err := <-chError
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// RunTaskPeriodically runs a task function periodically on the database object.
func (db *Database) RunTaskPeriodically(period time.Duration, task func(db *Database)) {
	chDone := make(chan struct{})
	for {
		go func() {
			defer func() { chDone <- struct{}{} }()
			task(db)
		}()
		time.Sleep(period)
		// makes sure task is finished before launching it again
		<-chDone
	}
}
