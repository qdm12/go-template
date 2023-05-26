package build

type Logger interface {
	Debugf(format string, args ...any)
	Error(s string)
}
