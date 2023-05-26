package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata"

	_ "github.com/breml/rootcerts"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qdm12/go-template/internal/config"
	"github.com/qdm12/go-template/internal/data"
	"github.com/qdm12/go-template/internal/health"
	"github.com/qdm12/go-template/internal/metrics"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/go-template/internal/processor"
	"github.com/qdm12/go-template/internal/server"
	"github.com/qdm12/golibs/params"
	"github.com/qdm12/goservices"
	"github.com/qdm12/goservices/httpserver"
	"github.com/qdm12/gosplash"
	"github.com/qdm12/log"
)

var (
	// Values set by the build system.
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	buildInfo := models.BuildInformation{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}

	args := os.Args

	logger := log.New()

	env := params.New()

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, env)
	}()

	select {
	case <-ctx.Done():
		logger.Warn("Caught OS signal, shutting down\n")
		stop()
	case err := <-errorCh:
		close(errorCh)
		if err == nil { // expected exit such as healthcheck query
			os.Exit(0)
		}
		logger.Error("Fatal error: " + err.Error())
		os.Exit(1)
	}

	err := <-errorCh
	if err != nil {
		logger.Error("shutdown error: " + err.Error())
	}
	os.Exit(1)
}

func _main(ctx context.Context, buildInfo models.BuildInformation,
	args []string, logger log.LoggerInterface, env params.Interface) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if health.IsClientMode(args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()

		var config config.Health
		_, err := config.Read(env)
		if err != nil {
			return err
		}

		return client.Query(ctx, config.Address)
	}

	announcementExpiration, err := time.Parse("2006-01-02", "2021-07-14")
	if err != nil {
		return err
	}
	splashLines := gosplash.MakeLines(gosplash.Settings{
		User:          "qdm12",
		Repository:    "go-template",
		Authors:       []string{"github.com/qdm12"},
		Emails:        []string{"quentin.mcgaw@gmail.com"},
		Version:       buildInfo.Version,
		Commit:        buildInfo.Commit,
		BuildDate:     buildInfo.BuildDate,
		Announcement:  "",
		AnnounceExp:   announcementExpiration,
		PaypalUser:    "qmcgaw",
		GithubSponsor: "qdm12",
	})
	fmt.Println(strings.Join(splashLines, "\n"))

	var config config.Config
	warnings, err := config.Read(env)
	for _, warning := range warnings {
		logger.Warn(warning)
	}
	if err != nil {
		return err
	}

	logger.Patch(log.SetLevel(config.Log.Level))

	db, err := setupDatabase(config.Store, logger)
	if err != nil {
		return err
	}

	proc := processor.NewProcessor(db)

	metricsServerSettings := httpserver.Settings{
		Name:    ptrTo("metrics"),
		Handler: promhttp.Handler(),
		Address: &config.Metrics.Address,
		Logger:  logger.New(log.SetComponent("metrics server")),
	}
	metricsServer, err := httpserver.New(metricsServerSettings)
	if err != nil {
		return fmt.Errorf("creating metrics server: %w", err)
	}
	const registerMetrics = true
	metrics, err := metrics.New(registerMetrics)
	if err != nil {
		return err
	}

	serverLogger := logger.New(log.SetComponent("http server"))
	serverSettings := httpserver.Settings{
		Name:    ptrTo("main"),
		Handler: server.NewRouter(config.HTTP, serverLogger, metrics, buildInfo, proc),
		Address: &config.HTTP.Address,
		Logger:  serverLogger,
	}
	mainServer, err := httpserver.New(serverSettings)
	if err != nil {
		return fmt.Errorf("creating main server: %w", err)
	}

	heathcheckLogger := logger.New(log.SetComponent("healthcheck"))
	healthcheck := func() error { return nil }
	healthServerHandler := health.NewHandler(heathcheckLogger, healthcheck)
	healthServerSettings := httpserver.Settings{
		Name:    ptrTo("health"),
		Handler: healthServerHandler,
		Address: &config.Health.Address,
		Logger:  heathcheckLogger,
	}
	healthServer, err := httpserver.New(healthServerSettings)
	if err != nil {
		return fmt.Errorf("creating health server: %w", err)
	}

	sequenceSettings := goservices.SequenceSettings{
		ServicesStart: []goservices.Service{db, metricsServer, healthServer, mainServer},
		ServicesStop:  []goservices.Service{mainServer, db, healthServer, metricsServer},
	}
	services, err := goservices.NewSequence(sequenceSettings)
	if err != nil {
		return fmt.Errorf("creating sequence of services: %w", err)
	}

	runError, err := services.Start()
	if err != nil {
		return fmt.Errorf("starting services: %w", err)
	}

	select {
	case <-ctx.Done():
		err = services.Stop()
		if err != nil {
			return fmt.Errorf("stopping services: %w", err)
		}
		return nil
	case err = <-runError:
		return fmt.Errorf("one service crashed, all services stopped: %w", err)
	}
}

var errDatabaseTypeUnknown = errors.New("database type is unknown")

func setupDatabase(c config.Store, logger log.LeveledLogger) (db data.Database, err error) {
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

func ptrTo[T any](x T) *T { return &x }
