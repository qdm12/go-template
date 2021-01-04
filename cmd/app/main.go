package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/qdm12/REPONAME_GITHUB/internal/data"
	"github.com/qdm12/REPONAME_GITHUB/internal/health"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
	"github.com/qdm12/REPONAME_GITHUB/internal/params"
	"github.com/qdm12/REPONAME_GITHUB/internal/processor"
	"github.com/qdm12/REPONAME_GITHUB/internal/server"
	"github.com/qdm12/REPONAME_GITHUB/internal/splash"
	"github.com/qdm12/golibs/crypto"
	"github.com/qdm12/golibs/logging"
)

var (
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

func main() {
	buildInfo := models.BuildInformation{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	paramsReader := params.NewReader()

	encoding, level, err := paramsReader.GetLoggerConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger, err := logging.NewLogger(encoding, level)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	args := os.Args

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, paramsReader)
	}()

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	select {
	case signal := <-signalsCh:
		logger.Warn("Caught OS signal %s, shutting down", signal)
	case err := <-errorCh:
		logger.Error(err)
		close(errorCh)
	}

	cancel()

	const shutdownGracePeriod = 5 * time.Second
	timer := time.NewTimer(shutdownGracePeriod)
	select {
	case <-errorCh:
		if !timer.Stop() {
			<-timer.C
		}
		logger.Info("Shutdown successful")
	case <-timer.C:
		logger.Warn("Shutdown timed out")
	}

	os.Exit(1)
}

func _main(ctx context.Context, buildInfo models.BuildInformation,
	_ []string, logger logging.Logger, paramsReader params.Reader) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if health.IsClientMode(os.Args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()
		if err := client.Query(ctx); err != nil {
			return err
		}
		return nil
	}

	fmt.Println(splash.Splash(buildInfo))
	listeningPort, warning, err := paramsReader.GetListeningPort()
	if len(warning) > 0 {
		logger.Warn(warning)
	}
	if err != nil {
		return err
	}
	rootURL, err := paramsReader.GetRootURL()
	if err != nil {
		return err
	}
	db, err := setupDatabase(paramsReader, logger)
	if err != nil {
		return err
	}
	defer db.Close()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	crypto := crypto.NewCrypto()
	proc := processor.NewProcessor(db, crypto)

	serverLogger := logger.WithPrefix("http server: ")
	address := fmt.Sprintf("0.0.0.0:%d", listeningPort)
	server := server.New(address, rootURL, serverLogger, buildInfo, proc)
	wg.Add(1)
	go server.Run(ctx, wg)

	const healthServerAddr = "127.0.0.1:9999"
	healthcheck := func() error { return nil }
	healthServer := health.NewServer(healthServerAddr, logger.WithPrefix("healthcheck server: "), healthcheck)
	wg.Add(1)
	go healthServer.Run(ctx, wg)

	<-ctx.Done()
	return nil
}

func setupDatabase(paramsReader params.Reader, logger logging.Logger) (db data.Database, err error) {
	databaseType := "memory"
	switch databaseType { // TODO env variable
	case "memory":
		return data.NewMemory()
	case "json":
		return data.NewJSON("data.json")
	case "postgres":
		dbHost, dbUser, dbPassword, dbName, err := paramsReader.GetDatabaseDetails()
		if err != nil {
			return nil, err
		}
		return data.NewPostgres(dbHost, dbUser, dbPassword, dbName, logger)
	default:
		return nil, fmt.Errorf("database type %q not supported", databaseType)
	}
}
