package config

import (
	"errors"
	"fmt"

	"github.com/qdm12/golibs/params"
)

type StoreType string

const (
	MemoryStoreType   StoreType = "memory"
	JSONStoreType     StoreType = "json"
	PostgresStoreType StoreType = "postgres"
)

type Store struct {
	Type     StoreType
	Memory   MemoryStore
	JSON     JSONStore
	Postgres Postgres
}

var (
	ErrMemoryStoreConfig = errors.New("memory config")
	ErrJSONStoreConfig   = errors.New("JSON config")
	ErrPostgresConfig    = errors.New("Postgres config")
)

func (s *Store) get(env params.Env) (warning string, err error) {
	s.Type, err = s.getType(env)
	if err != nil {
		return "", err
	}

	switch s.Type {
	case MemoryStoreType:
		err = s.Memory.get(env)
		if err != nil {
			err = fmt.Errorf("%w: %s", ErrMemoryStoreConfig, err)
		}
	case JSONStoreType:
		err = s.JSON.get(env)
		if err != nil {
			err = fmt.Errorf("%w: %s", ErrJSONStoreConfig, err)
		}
	case PostgresStoreType:
		warning, err = s.Postgres.get(env)
		if err != nil {
			err = fmt.Errorf("%w: %s", ErrPostgresConfig, err)
		}
	}

	return warning, err
}

func (s *Store) getType(env params.Env) (t StoreType, err error) {
	const envKey = "STORE_TYPE"
	possibilities := []string{"memory", "json", "postgres"}
	value, err := env.Inside(envKey, possibilities, params.Default("memory"))
	if err != nil {
		return t, err
	}
	return StoreType(value), nil
}
