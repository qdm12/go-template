package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/crypto"
	"github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/server"
	"github.com/qdm12/golibs/signals"

	"github.com/qdm12/REPONAME_GITHUB/internal/data"
	"github.com/qdm12/REPONAME_GITHUB/internal/env"
	"github.com/qdm12/REPONAME_GITHUB/internal/handlers"
	"github.com/qdm12/REPONAME_GITHUB/internal/params"
	"github.com/qdm12/REPONAME_GITHUB/internal/processor"
	"github.com/qdm12/REPONAME_GITHUB/internal/splash"
)

func main() {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel, -1)
	if err != nil {
		panic(err)
	}
	paramsGetter := params.NewGetter()
	encoding, level, nodeID, err := paramsGetter.GetLoggerConfig()
	if err != nil {
		logger.Error(err)
	} else {
		logger, err = logging.NewLogger(encoding, level, nodeID)
		if err != nil {
			panic(err)
		}
	}
	if healthcheck.Mode(os.Args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		if err := healthcheck.Query(); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	fmt.Println(splash.Splash(paramsGetter))
	e := env.NewEnv(logger)
	gotifyURL, err := paramsGetter.GetGotifyURL()
	e.FatalOnError(err)
	if gotifyURL != nil {
		gotifyToken, err := paramsGetter.GetGotifyToken()
		e.FatalOnError(err)
		e.SetGotify(admin.NewGotify(*gotifyURL, gotifyToken, &http.Client{Timeout: time.Second}))
	}
	listeningPort, warning, err := paramsGetter.GetListeningPort()
	e.FatalOnError(err)
	if len(warning) > 0 {
		logger.Warn(warning)
	}
	rootURL, err := paramsGetter.GetRootURL()
	e.FatalOnError(err)
	HTTPTimeout, err := paramsGetter.GetHTTPTimeout()
	e.FatalOnError(err)
	e.SetClient(network.NewClient(HTTPTimeout))

	var db data.Database
	switch "memory" { // TODO env variable
	case "memory":
		db, err = data.NewMemory()
		e.FatalOnError(err)
	case "json":
		db, err = data.NewJSON("data.json")
		e.FatalOnError(err)
	case "postgres":
		dbHost, dbUser, dbPassword, dbName, err := paramsGetter.GetDatabaseDetails()
		e.FatalOnError(err)
		db, err = data.NewPostgres(dbHost, dbUser, dbPassword, dbName, logger)
		e.FatalOnError(err)
	}
	e.SetDb(db)
	defer e.Shutdown()
	go signals.WaitForExit(e.ShutdownFromSignal)
	for _, err := range network.NewConnectivity(3 * time.Second).Checks("google.com") {
		e.Warn(err)
	}
	crypto := crypto.NewCrypto()
	proc := processor.NewProcessor(db, crypto)
	productionHandlerFunc := handlers.NewHandler(rootURL, proc, logger).GetHandlerFunc()
	healthcheckHandlerFunc := healthcheck.GetHandler(func() error { return nil })
	e.Notify(1, "About to launch HTTP servers")
	e.Info("About to launch HTTP servers")
	serverErrs := server.RunServers(
		server.Settings{Name: "production", Addr: "0.0.0.0:" + listeningPort, Handler: productionHandlerFunc},
		server.Settings{Name: "healthcheck", Addr: "127.0.0.1:9999", Handler: healthcheckHandlerFunc},
	)
	for _, err := range serverErrs {
		e.CheckError(err)
	}
	if len(serverErrs) > 0 {
		e.Fatal(serverErrs)
	}
}
