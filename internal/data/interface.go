package data

import (
	"context"

	"github.com/qdm12/REPONAME_GITHUB/internal/config"
	"github.com/qdm12/REPONAME_GITHUB/internal/data/json"
	"github.com/qdm12/REPONAME_GITHUB/internal/data/memory"
	"github.com/qdm12/REPONAME_GITHUB/internal/data/psql"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
	"github.com/qdm12/golibs/logging"
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

func NewPostgres(config config.Postgres, logger logging.Logger) (Database, error) {
	return psql.NewDatabase(config, logger)
}
