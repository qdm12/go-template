package env

import (
	"os"

	"github.com/qdm12/REPONAME_GITHUB/internal/data"

	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

// Env contains objects and methods necessary to the main function.
// These are created at start and are needed to the top-level
// working management of the program.
type Env interface {
	SetGotify(gotify admin.Gotify)
	SetDb(db data.Database)
	SetClient(client network.Client)
	Notify(priority int, args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	CheckError(err error)
	FatalOnError(err error)
	ShutdownFromSignal(signal string) (exitCode int)
	Fatal(args ...interface{})
	Shutdown() (exitCode int)
}

type env struct {
	client network.Client
	logger logging.Logger
	gotify admin.Gotify
	db     data.Database
}

// NewEnv creates a new Env object
func NewEnv(logger logging.Logger) Env {
	return &env{logger: logger}
}

func (e *env) SetGotify(gotify admin.Gotify) {
	e.gotify = gotify
}

func (e *env) SetDb(db data.Database) {
	e.db = db
}

func (e *env) SetClient(client network.Client) {
	e.client = client
}

// Notify sends a notification to the Gotify server.
func (e *env) Notify(priority int, args ...interface{}) {
	if e.gotify == nil {
		return
	}
	if err := e.gotify.Notify("Program name", priority, args...); err != nil {
		e.logger.Error(err)
	}
}

// Info logs a message and sends a notification to the Gotify server.
func (e *env) Info(args ...interface{}) {
	e.logger.Info(args...)
	e.Notify(1, args...)
}

// Warn logs a message and sends a notification to the Gotify server.
func (e *env) Warn(args ...interface{}) {
	e.logger.Warn(args...)
	e.Notify(2, args...)
}

// CheckError logs an error and sends a notification to the Gotify server
// if the error is not nil.
func (e *env) CheckError(err error) {
	if err == nil {
		return
	}
	s := err.Error()
	e.logger.Error(s)
	if len(s) > 100 {
		s = s[:100] + "..." // trim down message for notification
	}
	e.Notify(3, s)
}

// FatalOnError calls Fatal if the error is not nil.
func (e *env) FatalOnError(err error) {
	if err != nil {
		e.Fatal(err)
	}
}

// Shutdown cleanly exits the program by closing all connections,
// databases and syncing the loggers.
func (e *env) Shutdown() (exitCode int) {
	defer func() {
		if err := e.logger.Sync(); err != nil {
			exitCode = 99
		}
	}()
	if e.client != nil {
		e.client.Close()
	}
	if e.db != nil {
		if err := e.db.Close(); err != nil {
			e.logger.Error(err)
			return 1
		}
	}
	return 0
}

// ShutdownFromSignal logs a warning, sends a notification to Gotify and shutdowns
// the program cleanly when a OS level signal is received. It should be passed as a
// callback to a function which would catch such signal.
func (e *env) ShutdownFromSignal(signal string) (exitCode int) {
	e.logger.Warn("Program stopped with signal %s", signal)
	e.Notify(1, "Caught OS signal "+signal)
	return e.Shutdown()
}

// Fatal logs an error, sends a notification to Gotify and shutdowns the program.
// It exits the program with an exit code of 1.
func (e *env) Fatal(args ...interface{}) {
	e.logger.Error(args...)
	e.Notify(4, args...)
	_ = e.Shutdown()
	os.Exit(1)
}
