package env

import (
	"github.com/qdm12/go-template/internal/config/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func readDatabase() (database settings.Database) {
	database.Type = env.StringPtr("STORE_TYPE")
	database.Memory = readMemoryDatabase()
	database.JSON = readJSONDatabase()
	database.Postgres = readPostgresDatabase()
	return database
}

func readMemoryDatabase() (
	database settings.MemoryDatabase) {
	return database
}

func readJSONDatabase() (
	database settings.JSONDatabase) {
	database.Filepath = env.Get("JSON_FILEPATH", env.ForceLowercase(false))
	return database
}

func readPostgresDatabase() (
	database settings.PostgresDatabase) {
	database.Address = env.Get("POSTGRES_ADDRESS")
	database.User = env.Get("POSTGRES_USER", env.ForceLowercase(false))
	database.Password = env.Get("POSTGRES_PASSWORD", env.ForceLowercase(false))
	database.Database = env.Get("POSTGRES_DATABASE", env.ForceLowercase(false))
	return database
}
