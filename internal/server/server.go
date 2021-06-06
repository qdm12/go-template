// Package server implements an HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/go-template/internal/config"
	"github.com/qdm12/go-template/internal/metrics"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/go-template/internal/processor"
	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context) error
}

type server struct {
	address string
	logger  logging.Logger
	handler http.Handler
}

func New(c config.HTTP, proc processor.Processor,
	logger logging.Logger, metrics metrics.Metrics,
	buildInfo models.BuildInformation) Server {
	handler := newRouter(c, logger, metrics, buildInfo, proc)
	return &server{
		address: c.Address,
		logger:  logger,
		handler: handler,
	}
}

var (
	ErrCrashed  = errors.New("server crashed")
	ErrShutdown = errors.New("server could not be shutdown")
)

func (s *server) Run(ctx context.Context) error {
	server := http.Server{Addr: s.address, Handler: s.handler}
	shutdownErrCh := make(chan error)
	go func() {
		<-ctx.Done()
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		shutdownErrCh <- server.Shutdown(shutdownCtx)
	}()

	s.logger.Info("listening on " + s.address)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(ctx.Err(), context.Canceled) { // server crashed
		return fmt.Errorf("%w: %s", ErrCrashed, err)
	}

	if err := <-shutdownErrCh; err != nil {
		return fmt.Errorf("%w: %s", ErrShutdown, err)
	}

	return nil
}
