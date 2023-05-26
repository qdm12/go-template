package build

type Logger interface {
	Debugf(format string, args ...interface{})
	Error(s string)
}
