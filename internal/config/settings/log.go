package settings

import (
	"errors"
	"fmt"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
	"github.com/qdm12/log"
)

type Log struct {
	Level string
}

func (l *Log) setDefaults() {
	l.Level = gosettings.DefaultComparable(l.Level, log.LevelInfo.String())
}

var (
	ErrLogLevelUnknown = errors.New("log level is unknown")
)

func (l *Log) validate() (err error) {
	_, err = log.ParseLevel(l.Level)
	if err != nil {
		return fmt.Errorf("log level: %w", err)
	}
	return nil
}

func (l *Log) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Log settings:")
	node.Appendf("Level: %s", l.Level)
	return node
}

func (l *Log) copy() (copied Log) {
	return Log{
		Level: l.Level,
	}
}

func (l *Log) overrideWith(other Log) {
	l.Level = gosettings.OverrideWithComparable(l.Level, other.Level)
}

func (l *Log) read(r *reader.Reader) {
	l.Level = r.String("LOG_LEVEL")
}
