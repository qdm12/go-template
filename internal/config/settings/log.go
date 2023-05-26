package settings

import (
	"errors"
	"fmt"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
	"github.com/qdm12/log"
)

type Log struct {
	Level *log.Level
}

func (l *Log) setDefaults() {
	l.Level = gosettings.DefaultPointer(l.Level, log.LevelInfo)
}

var (
	ErrLogLevelUnknown = errors.New("log level is unknown")
)

func (l *Log) validate() (err error) {
	switch *l.Level {
	case log.LevelDebug, log.LevelInfo, log.LevelWarn, log.LevelError:
	default:
		return fmt.Errorf("%w: %d", ErrLogLevelUnknown, *l.Level)
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
		Level: gosettings.CopyPointer(l.Level),
	}
}

func (l *Log) mergeWith(other Log) {
	l.Level = gosettings.MergeWithPointer(l.Level, other.Level)
}

func (l *Log) overrideWith(other Log) {
	l.Level = gosettings.OverrideWithPointer(l.Level, other.Level)
}
