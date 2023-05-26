package server

type Logger interface {
	Infof(format string, args ...interface{})
	Error(s string)
}
