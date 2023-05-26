// Package psql implements a data store using a client to a
// Postgres database.
package psql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/qdm12/go-template/internal/config/settings"
	"github.com/qdm12/goservices"
)

// Database is the Postgres implementation of the database store.
type Database struct {
	startStopMutex sync.Mutex
	running        bool
	sql            *sql.DB
	logger         Logger
}

// NewDatabase creates a database connection pool in DB and pings the database.
func NewDatabase(config settings.PostgresDatabase, logger Logger) (*Database, error) {
	connStr := "postgres://" + config.User + ":" + config.Password +
		"@" + config.Address + "/" + config.Address + "?sslmode=disable&connect_timeout=1"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Database{
		sql:    db,
		logger: logger,
	}, nil
}

func (db *Database) String() string {
	return "postgres database"
}

// Start pings the database, and if it fails, retries up to 3 times
// before returning a start error.
func (db *Database) Start() (runError <-chan error, err error) {
	db.startStopMutex.Lock()
	defer db.startStopMutex.Unlock()

	if db.running {
		return nil, fmt.Errorf("%w", goservices.ErrAlreadyStarted)
	}

	fails := 0
	const maxFails = 3
	const sleepDuration = 200 * time.Millisecond
	var totalTryTime time.Duration
	for {
		err = db.sql.Ping()
		if err == nil {
			break
		}
		fails++
		if fails == maxFails {
			return nil, fmt.Errorf("failed connecting to database after %d tries in %s: %w", fails, totalTryTime, err)
		}
		time.Sleep(sleepDuration)
		totalTryTime += sleepDuration
	}

	db.running = true
	// TODO have periodic ping to check connection is still alive
	// and signal through the run error channel.
	return nil, nil
}

// Stop stops the database and closes the connection.
func (db *Database) Stop() (err error) {
	db.startStopMutex.Lock()
	defer db.startStopMutex.Unlock()
	if !db.running {
		return fmt.Errorf("%w", goservices.ErrAlreadyStopped)
	}

	err = db.sql.Close()
	if err != nil {
		return fmt.Errorf("closing database connection: %w", err)
	}

	db.running = false
	return nil
}
