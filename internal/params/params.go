package params

import (
	"net/url"
	"time"

	"github.com/qdm12/golibs/logging"
	libparams "github.com/qdm12/golibs/params"
)

type Getter interface {
	GetListeningPort() (listeningPort, warning string, err error)
	GetLoggerConfig() (encoding logging.Encoding, level logging.Level, nodeID int, err error)
	GetGotifyURL(setters ...libparams.GetEnvSetter) (*url.URL, error)
	GetGotifyToken(setters ...libparams.GetEnvSetter) (token string, err error)
	GetRootURL(setters ...libparams.GetEnvSetter) (rootURL string, err error)
	GetHTTPTimeout() (duration time.Duration, err error)
	GetDatabaseDetails() (hostname, user, password, dbName string, err error)

	// Version getters
	GetVersion() string
	GetBuildDate() string
	GetVcsRef() string
}

type getter struct {
	envParams libparams.EnvParams
}

func NewGetter() Getter {
	return &getter{
		envParams: libparams.NewEnvParams(),
	}
}

// GetDataDir obtains the data directory from the environment
// variable DATADIR
func (g *getter) GetDataDir(currentDir string) (string, error) {
	return g.envParams.GetEnv("DATADIR", libparams.Default(currentDir+"/data"))
}

func (g *getter) GetListeningPort() (listeningPort, warning string, err error) {
	return g.envParams.GetListeningPort()
}

func (g *getter) GetLoggerConfig() (encoding logging.Encoding, level logging.Level, nodeID int, err error) {
	return g.envParams.GetLoggerConfig()
}

func (g *getter) GetGotifyURL(setters ...libparams.GetEnvSetter) (*url.URL, error) {
	return g.envParams.GetGotifyURL()
}

func (g *getter) GetGotifyToken(setters ...libparams.GetEnvSetter) (token string, err error) {
	return g.envParams.GetGotifyToken()
}

func (g *getter) GetRootURL(setters ...libparams.GetEnvSetter) (rootURL string, err error) {
	return g.envParams.GetRootURL()
}

func (g *getter) GetExeDir() (dir string, err error) {
	return g.envParams.GetExeDir()
}

func (g *getter) GetHTTPTimeout() (duration time.Duration, err error) {
	return g.envParams.GetHTTPTimeout(libparams.Default("10s"))
}

func (g *getter) GetDatabaseDetails() (hostname, user, password, dbName string, err error) {
	return g.envParams.GetDatabaseDetails()
}

func (g *getter) GetVersion() string {
	version, _ := g.envParams.GetEnv("VERSION", libparams.Default("?"), libparams.CaseSensitiveValue())
	return version
}

func (g *getter) GetBuildDate() string {
	buildDate, _ := g.envParams.GetEnv("BUILD_DATE", libparams.Default("?"), libparams.CaseSensitiveValue())
	return buildDate
}

func (g *getter) GetVcsRef() string {
	buildDate, _ := g.envParams.GetEnv("VCS_REF", libparams.Default("?"), libparams.CaseSensitiveValue())
	return buildDate
}
