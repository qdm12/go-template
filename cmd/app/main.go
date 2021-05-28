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
	"github.com/qdm12/go-template/internal/config"
	"github.com/qdm12/go-template/internal/data"
	"github.com/qdm12/go-template/internal/health"
	"github.com/qdm12/go-template/internal/metrics"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/go-template/internal/processor"
	"github.com/qdm12/go-template/internal/server"
	"github.com/qdm12/go-template/internal/splash"
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
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	configReader := config.NewReader()

	logger := logging.NewParent(logging.Settings{})

	args := os.Args

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, configReader)
	}()

	select {
	case <-ctx.Done():
		logger.Warn("Caught OS signal, shutting down\n")
		stop()
	case err := <-errorCh:
		close(errorCh)
		if err == nil { // expected exit such as healthcheck
			os.Exit(0)
		}
		logger.Error(err)
	}

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
	args []string, logger logging.ParentLogger, configReader config.Reader) error {
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

	logger = logger.NewChild(logging.Settings{Level: config.Log.Level})

	db, err := setupDatabase(config.Store, logger)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	crashed := make(chan error)

	proc := processor.NewProcessor(db)

	metricsLogger := logger.NewChild(logging.Settings{Prefix: "metrics server: "})
	metricsServer := metrics.NewServer(config.Metrics.Address, metricsLogger)
	const registerMetrics = true
	metrics, err := metrics.New(registerMetrics)
	if err != nil {
		return err
	}
	wg.Add(1)
	go metricsServer.Run(ctx, wg, crashed)

	serverLogger := logger.NewChild(logging.Settings{Prefix: "http server: "})
	server := server.New(config.HTTP, proc, serverLogger, metrics, buildInfo)
	wg.Add(1)
	go server.Run(ctx, wg, crashed)

	healthcheck := func() error { return nil }
	heathcheckLogger := logger.NewChild(logging.Settings{Prefix: "healthcheck: "})
	healthServer := health.NewServer(config.Health.Address, heathcheckLogger, healthcheck)
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
