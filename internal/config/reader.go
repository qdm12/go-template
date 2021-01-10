package config

import (
	"github.com/qdm12/golibs/params"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Reader

type Reader interface {
	ReadConfig() (c Config, warnings []string, err error)
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
