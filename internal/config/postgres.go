package config

import (
	"github.com/qdm12/golibs/params"
)

type Postgres struct {
	Address  string
	User     string
	Password string
	Database string
}

func (p *Postgres) get(env params.Env) (warning string, err error) {
	p.User, err = p.getUser(env)
	if err != nil {
		return "", err
	}

	p.Password, err = p.getPassword(env)
	if err != nil {
		return "", err
	}

	p.Database, err = p.getDatabase(env)
	if err != nil {
		return "", err
	}

	p.Address, warning, err = p.getAddress(env)
	if err != nil {
		return warning, err
	}

	return warning, nil
}

func (p *Postgres) getAddress(env params.Env) (address, warning string, err error) {
	const envKey = "STORE_POSTGRES_ADDRESS"
	options := []params.OptionSetter{
		params.Default("psql:5432"),
	}
	return env.ListeningAddress(envKey, options...)
}

func (p *Postgres) getUser(env params.Env) (user string, err error) {
	const envKey = "STORE_POSTGRES_USER"
	options := []params.OptionSetter{
		params.Default("postgres"),
		params.CaseSensitiveValue(),
		params.Unset(),
	}
	return env.Get(envKey, options...)
}

func (p *Postgres) getPassword(env params.Env) (password string, err error) {
	const envKey = "STORE_POSTGRES_PASSWORD"
	options := []params.OptionSetter{
		params.Default("postgres"),
		params.CaseSensitiveValue(),
		params.Unset(),
	}
	return env.Get(envKey, options...)
}

func (p *Postgres) getDatabase(env params.Env) (database string, err error) {
	const envKey = "STORE_POSTGRES_DATABASE"
	options := []params.OptionSetter{
		params.Default("database"),
	}
	return env.Get(envKey, options...)
}
