package params

import (
	"github.com/qdm12/golibs/logging"
	libparams "github.com/qdm12/golibs/params"
)

type Reader interface {
	GetListeningPort() (listeningPort uint16, warning string, err error)
	GetLoggerConfig() (encoding logging.Encoding, level logging.Level, err error)
	GetRootURL(setters ...libparams.GetEnvSetter) (rootURL string, err error)
	GetDatabaseDetails() (hostname, user, password, dbName string, err error)
}

type reader struct {
	envParams libparams.EnvParams
}

func NewReader() Reader {
	return &reader{
		envParams: libparams.NewEnvParams(),
	}
}

func (r *reader) GetListeningPort() (listeningPort uint16, warning string, err error) {
	return r.envParams.GetListeningPort("LISTENING_PORT")
}

func (r *reader) GetLoggerConfig() (encoding logging.Encoding, level logging.Level, err error) {
	return r.envParams.GetLoggerConfig()
}

func (r *reader) GetRootURL(setters ...libparams.GetEnvSetter) (rootURL string, err error) {
	return r.envParams.GetRootURL()
}

func (r *reader) GetDatabaseDetails() (hostname, user, password, dbName string, err error) {
	return r.envParams.GetDatabaseDetails()
}
