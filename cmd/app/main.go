package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	_ "github.com/lib/pq"
	"github.com/qdm12/go-template/internal/config"
	"github.com/qdm12/go-template/internal/data"
	"github.com/qdm12/go-template/internal/health"
	"github.com/qdm12/go-template/internal/metrics"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/go-template/internal/processor"
	"github.com/qdm12/go-template/internal/server"
	"github.com/qdm12/go-template/internal/shutdown"
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

		config, _, err := configReader.ReadConfig()
		if err != nil {
			return err
		}

		return client.Query(ctx, config.Health.Address)
	}

	fmt.Println(splash.Splash(buildInfo))

	config, warnings, err := configReader.ReadConfig()
	for _, warning := range warnings {
		logger.Warn(warning)
	}
	if err != nil {
		return err
	}

	shutdownServersGroup := shutdown.NewGroup("servers: ")
	shutdownStoreGroup := shutdown.NewGroup("store: ")

	logger = logger.NewChild(logging.Settings{Level: config.Log.Level})

	db, err := setupDatabase(config.Store, logger)
	if err != nil {
		return err
	}

	proc := processor.NewProcessor(db)

	metricsLogger := logger.NewChild(logging.Settings{Prefix: "metrics server: "})
	metricsServer := metrics.NewServer(config.Metrics.Address, metricsLogger)
	const registerMetrics = true
	metrics, err := metrics.New(registerMetrics)
	if err != nil {
		return err
	}
	metricsServerCtx, metricsServerDone := shutdownServersGroup.Add("metrics", time.Second)
	go func() {
		defer close(metricsServerDone)
		if err := metricsServer.Run(metricsServerCtx); err != nil {
			logger.Error(err.Error())
		}
	}()

	serverLogger := logger.NewChild(logging.Settings{Prefix: "http server: "})
	mainServer := server.New(config.HTTP, proc, serverLogger, metrics, buildInfo)
	serverCtx, serverDone := shutdownServersGroup.Add("server", time.Second)
	go func() {
		defer close(serverDone)
		if err := mainServer.Run(serverCtx); err != nil {
			logger.Error(err.Error())
			if errors.Is(err, server.ErrCrashed) {
				cancel() // stop other routines
			}
		}
	}()

	healthcheck := func() error { return nil }
	heathcheckLogger := logger.NewChild(logging.Settings{Prefix: "healthcheck: "})
	healthServer := health.NewServer(config.Health.Address, heathcheckLogger, healthcheck)
	healthServerCtx, healthServerDone := shutdownServersGroup.Add("health", time.Second)
	go func() {
		defer close(healthServerDone)
		if err := healthServer.Run(healthServerCtx); err != nil {
			logger.Error(err.Error())
		}
	}()

	// Adapt db.Close to the shutdown logic
	dbCloseCtx, dbCloseDone := shutdownStoreGroup.Add("close", time.Second)
	go func() {
		<-dbCloseCtx.Done()
		db.Close()
		close(dbCloseDone)
	}()

	shutdownOrder := shutdown.NewOrder()
	shutdownOrder.Append(shutdownServersGroup, shutdownStoreGroup)

	<-ctx.Done()
	return shutdownOrder.Shutdown(time.Second, logger)
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
