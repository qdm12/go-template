// Package data contains a Database interface with multiple implementations.
package data

import (
	"github.com/qdm12/go-template/internal/config/settings"
	"github.com/qdm12/go-template/internal/data/json"
	"github.com/qdm12/go-template/internal/data/memory"
	"github.com/qdm12/go-template/internal/data/psql"
	"github.com/qdm12/log"
)

func NewMemory() (db *memory.Database, err error) {
	return memory.NewDatabase()
}

func NewJSON(filepath string) (db *json.Database, err error) {
	memoryDatabase, err := memory.NewDatabase()
	if err != nil {
		return nil, err
	}
	return json.NewDatabase(memoryDatabase, filepath), nil
}

func NewPostgres(config settings.PostgresDatabase, logger log.LeveledLogger) (
	db *psql.Database, err error,
) {
	return psql.NewDatabase(config, logger)
}
