package log

type Logger interface {
	Infof(format string, args ...any)
}
