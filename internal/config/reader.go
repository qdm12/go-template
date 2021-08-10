package config

import (
	"github.com/qdm12/golibs/params"
)

var _ ReaderInterface = (*Reader)(nil)

type ReaderInterface interface {
	ConfReader
	HealthReader
}

type Reader struct {
	env params.Env
}

func NewReader() *Reader {
	return &Reader{
		env: params.NewEnv(),
	}
}

type ConfReader interface {
	ReadConfig() (c Config, warnings []string, err error)
}

// ReadConfig reads all the configuration and returns it.
func (r *Reader) ReadConfig() (c Config, warnings []string, err error) {
	warnings, err = c.get(r.env)
	return c, warnings, err
}

type HealthReader interface {
	ReadHealth() (h Health, err error)
}

// ReadHealth is used for the healthcheck query only.
func (r *Reader) ReadHealth() (h Health, err error) {
	// warning is ignored when reading in healthcheck client query mode.
	_, err = h.get(r.env)
	return h, err
}
