package config

import (
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
)

type Log struct {
	Encoding logging.Encoding
	Level    logging.Level
}

func (l *Log) get(env params.Env) (err error) {
	l.Encoding, err = l.getEncoding(env)
	if err != nil {
		return err
	}
	l.Level, err = l.getLevel(env)
	if err != nil {
		return err
	}
	return nil
}

func (l *Log) getEncoding(env params.Env) (encoding logging.Encoding, err error) {
	const envKey = "LOG_ENCODING"
	options := []params.OptionSetter{
		params.Default("console"),
	}
	return env.LogEncoding(envKey, options...)
}

func (l *Log) getLevel(env params.Env) (level logging.Level, err error) {
	const envKey = "LOG_LEVEL"
	options := []params.OptionSetter{
		params.Default("info"),
	}
	return env.LogLevel(envKey, options...)
}
