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
	"github.com/qdm12/go-template/internal/config/settings"
	"github.com/qdm12/go-template/internal/data"
	"github.com/qdm12/go-template/internal/health"
	"github.com/qdm12/go-template/internal/metrics"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/go-template/internal/processor"
	"github.com/qdm12/go-template/internal/server"
	"github.com/qdm12/go-template/internal/server/websocket"
	"github.com/qdm12/goservices"
	"github.com/qdm12/goservices/hooks"
	"github.com/qdm12/goservices/httpserver"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/reader/sources/env"
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
	buildInfo := models.BuildInformation{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}

	background := context.Background()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(background)

	args := os.Args

	logger := log.New()

	reader := reader.New(reader.Settings{
		Sources: []reader.Source{env.New(env.Settings{})},
		HandleDeprecatedKey: func(source string, deprecatedKey string, currentKey string) {
			logger.Warnf("Deprecated %s %q, please use %q instead", source, deprecatedKey, currentKey)
		},
	})

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, reader)
	}()

	// Wait for OS signal or run error
	var runError error
	select {
	case receivedSignal := <-signalCh:
		signal.Stop(signalCh)
		fmt.Println("")
		logger.Warn("Caught OS signal " + receivedSignal.String() + ", shutting down")
		cancel()
	case runError = <-errorCh:
		close(errorCh)
		if runError == nil { // expected exit such as healthcheck
			os.Exit(0)
		}
		logger.Error(runError.Error())
		cancel()
	}

	// Shutdown timed sequence, and force exit on second OS signal
	const shutdownGracePeriod = 5 * time.Second
	timer := time.NewTimer(shutdownGracePeriod)
	select {
	case shutdownErr := <-errorCh:
		timer.Stop()
		if shutdownErr != nil {
			logger.Warnf("Shutdown failed: %s", shutdownErr)
			os.Exit(1)
		}

		logger.Info("Shutdown successful")
		if runError != nil {
			os.Exit(1)
		}
		os.Exit(0)
	case <-timer.C:
		logger.Warn("Shutdown timed out")
		os.Exit(1)
	}
}

//nolint:cyclop
func _main(ctx context.Context, buildInfo models.BuildInformation,
	args []string, logger log.LoggerInterface, configReader *reader.Reader,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if health.IsClientMode(args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		var healthConfig settings.Health
		healthConfig.Read(configReader)
		healthConfig.SetDefaults()
		err := healthConfig.Validate()
		if err != nil {
			return fmt.Errorf("health configuration is invalid: %w", err)
		}

		client := health.NewClient()
		// TODO write listening address to file for the healthcheck to read
		// since the user can pass '' to listen on any available port.
		return client.Query(ctx, healthConfig.Address)
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

	var config settings.Settings
	err = config.Read(configReader)
	if err != nil {
		return fmt.Errorf("reading configuration: %w", err)
	}
	config.SetDefaults()
	err = config.Validate()
	if err != nil {
		return fmt.Errorf("configuration is invalid: %w", err)
	}

	logLevel, _ := log.ParseLevel(config.Log.Level) // level string already validated
	logger.Patch(log.SetLevel(logLevel))

	logger.Info(config.String())

	db, err := setupDatabase(config.Database, logger)
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
		Address: ptrTo(*config.HTTP.Address),
		Logger:  serverLogger,
	}
	mainServer, err := httpserver.New(serverSettings)
	if err != nil {
		return fmt.Errorf("creating main server: %w", err)
	}

	websocketServerLogger := logger.New(log.SetComponent("websocket server"))
	websocketServerSettings := httpserver.Settings{
		Name:    ptrTo("websocket"),
		Handler: websocket.New(),
		Address: config.HTTP.Websocket.Address,
		Logger:  websocketServerLogger,
	}
	websocketServer, err := httpserver.New(websocketServerSettings)
	if err != nil {
		return fmt.Errorf("creating websocket server: %w", err)
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

	servicesLogger := logger.New(log.SetComponent("services"))
	sequenceSettings := goservices.SequenceSettings{
		ServicesStart: []goservices.Service{db, metricsServer, healthServer, websocketServer, mainServer},
		ServicesStop:  []goservices.Service{mainServer, websocketServer, db, healthServer, metricsServer},
		Hooks:         hooks.NewWithLog(servicesLogger),
	}
	services, err := goservices.NewSequence(sequenceSettings)
	if err != nil {
		return fmt.Errorf("creating sequence of services: %w", err)
	}

	runError, err := services.Start(ctx)
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

type Database interface {
	String() string
	Start(ctx context.Context) (runError <-chan error, err error)
	Stop() (err error)
	CreateUser(ctx context.Context, user models.User) (err error)
	GetUserByID(ctx context.Context, id uint64) (user models.User, err error)
}

func setupDatabase(databaseSettings settings.Database, logger log.LeveledLogger) ( //nolint:ireturn
	db Database, err error,
) {
	switch *databaseSettings.Type {
	case settings.MemoryStoreType:
		return data.NewMemory()
	case settings.JSONStoreType:
		return data.NewJSON(databaseSettings.JSON.Filepath)
	case settings.PostgresStoreType:
		return data.NewPostgres(databaseSettings.Postgres, logger)
	default:
		return nil, fmt.Errorf("%w: %s", errDatabaseTypeUnknown, *databaseSettings.Type)
	}
}

func ptrTo[T any](x T) *T { return &x }

type ConfigSource interface {
	Read() (settings settings.Settings, err error)
	ReadHealth() (health settings.Health)
	String() string
}
