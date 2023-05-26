package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/golibs/params"
	"github.com/qdm12/log"
)

type Log struct {
	Level log.Level
}

func (l *Log) get(env params.Interface) (err error) {
	l.Level, err = readLogLevel()
	if err != nil {
		return err
	}
	return nil
}

func readLogLevel() (level log.Level, err error) {
	s := getCleanedEnv("LOG_LEVEL")
	if s == "" {
		return log.LevelInfo, nil //nolint:nilnil
	}

	level, err = parseLogLevel(s)
	if err != nil {
		return level, fmt.Errorf("environment variable LOG_LEVEL: %w", err)
	}

	return level, nil
}

var ErrLogLevelUnknown = errors.New("log level is unknown")

func parseLogLevel(s string) (level log.Level, err error) {
	switch strings.ToLower(s) {
	case "debug":
		return log.LevelDebug, nil
	case "info":
		return log.LevelInfo, nil
	case "warning":
		return log.LevelWarn, nil
	case "error":
		return log.LevelError, nil
	default:
		return level, fmt.Errorf(
			"%w: %q is not valid and can be one of debug, info, warning or error",
			ErrLogLevelUnknown, s)
	}
}
