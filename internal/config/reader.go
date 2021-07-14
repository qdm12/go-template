package config

import (
	"github.com/qdm12/golibs/params"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Reader

type Reader interface {
	// ReadConfig reads all the configuration and returns it.
	ReadConfig() (c Config, warnings []string, err error)
	// ReadHealth is used for the healthcheck query only.
	ReadHealth() (h Health, err error)
}

type reader struct {
	env params.Env
}

func NewReader() Reader {
	return &reader{
		env: params.NewEnv(),
	}
}

func (r *reader) ReadConfig() (c Config, warnings []string, err error) {
	warnings, err = c.get(r.env)
	return c, warnings, err
}

func (r *reader) ReadHealth() (h Health, err error) {
	// warning is ignored when reading in healthcheck client query mode.
	_, err = h.get(r.env)
	return h, err
}
