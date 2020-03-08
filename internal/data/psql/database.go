package psql

import (
	"database/sql"
	"time"

	"github.com/qdm12/golibs/crypto/random"
	"github.com/qdm12/golibs/logging"
)

// Database is the Postgres implementation of the database store
type Database struct {
	sql    *sql.DB
	logger logging.Logger
	random random.Random
}

// NewDatabase creates a database connection pool in DB and pings the database
func NewDatabase(host, user, password, database string, logger logging.Logger) (*Database, error) {
	connStr := "postgres://" + user + ":" + password + "@" + host + "/" + database + "?sslmode=disable&connect_timeout=1"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	fails := 0
	for {
		err = db.Ping()
		if err == nil {
			break
		}
		fails++
		if fails == 3 {
			return nil, err
		}
		time.Sleep(200 * time.Millisecond)
	}
	return &Database{db, logger, random.NewRandom()}, nil
}

// Close closes the database and prevents new queries from starting.
func (db *Database) Close() error {
	return db.sql.Close()
}

// PeriodicHealthcheck pings the database periodically and runs onFailure if an error occurs
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

// RunTaskPeriodically runs a task function periodically on the database object
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
