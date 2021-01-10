package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/qdm12/REPONAME_GITHUB/internal/config"
	"github.com/qdm12/REPONAME_GITHUB/internal/data"
	"github.com/qdm12/REPONAME_GITHUB/internal/health"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
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

	configReader := config.NewReader()

	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	args := os.Args

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, configReader)
	}()

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	select {
	case signal := <-signalsCh:
		logger.Warn("Caught OS signal %s, shutting down\n", signal)
	case err := <-errorCh:
		close(errorCh)
		if err == nil { // expected exit such as healthcheck
			os.Exit(0)
		}
		logger.Error(err)
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
	args []string, logger logging.Logger, configReader config.Reader) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if health.IsClientMode(args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()
		return client.Query(ctx)
	}

	fmt.Println(splash.Splash(buildInfo))

	config, warnings, err := configReader.ReadConfig()
	for _, warning := range warnings {
		logger.Warn(warning)
	}
	if err != nil {
		return err
	}

	logger, err = logging.NewLogger(config.Log.Encoding, config.Log.Level)
	if err != nil {
		return err
	}

	db, err := setupDatabase(config.Store, logger)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	crypto := crypto.NewCrypto()
	proc := processor.NewProcessor(db, crypto)

	serverLogger := logger.WithPrefix("http server: ")
	server := server.New(config.HTTP, serverLogger, buildInfo, proc)
	wg.Add(1)
	crashed := make(chan error)
	go server.Run(ctx, wg, crashed)

	healthcheck := func() error { return nil }
	healthServer := health.NewServer(
		config.Health.Address, logger.WithPrefix("healthcheck: "),
		healthcheck)
	wg.Add(1)
	go healthServer.Run(ctx, wg)

	select {
	case <-ctx.Done():
		wg.Wait()
		return db.Close()
	case err := <-crashed:
		cancel()
		wg.Wait()
		_ = db.Close()
		return err
	}
}

var errDatabaseTypeUnknown = errors.New("database type is unknown")

func setupDatabase(c config.Store, logger logging.Logger) (db data.Database, err error) {
	switch c.Type {
	case config.MemoryStoreType:
		return data.NewMemory()
	case config.JSONStoreType:
		return data.NewJSON(c.JSON.Filepath)
	case config.PostgresStoreType:
		return data.NewPostgres(c.Postgres, logger)
	default:
		return nil, fmt.Errorf("%w: %s", errDatabaseTypeUnknown, c.Type)
	}
}
