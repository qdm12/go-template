// Package data contains a Database interface with multiple implementations.
package data

import (
	"context"

	"github.com/qdm12/go-template/internal/config"
	"github.com/qdm12/go-template/internal/data/json"
	"github.com/qdm12/go-template/internal/data/memory"
	"github.com/qdm12/go-template/internal/data/psql"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/log"
)

type Database interface {
	Close() error
	CreateUser(ctx context.Context, user models.User) (err error)
	GetUserByID(ctx context.Context, id uint64) (user models.User, err error)
}

func NewMemory() (Database, error) {
	return memory.NewDatabase()
}

func NewJSON(filepath string) (Database, error) {
	memoryDatabase, err := memory.NewDatabase()
	if err != nil {
		return nil, err
	}
	return json.NewDatabase(memoryDatabase, filepath)
}

func NewPostgres(config config.Postgres, logger log.LeveledLogger) (Database, error) {
	return psql.NewDatabase(config, logger)
}
