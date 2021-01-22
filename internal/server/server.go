package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/REPONAME_GITHUB/internal/config"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
	"github.com/qdm12/REPONAME_GITHUB/internal/processor"
	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup, crashed chan<- error)
}

type server struct {
	address string
	logger  logging.Logger
	handler http.Handler
}

func New(c config.HTTP, logger logging.Logger,
	buildInfo models.BuildInformation, proc processor.Processor) Server {
	handler := newHandler(c.RootURL, logger, buildInfo, proc)
	return &server{
		address: c.Address,
		logger:  logger,
		handler: handler,
	}
}

func (s *server) Run(ctx context.Context, wg *sync.WaitGroup, crashed chan<- error) {
	defer wg.Done()
	server := http.Server{Addr: s.address, Handler: s.handler}
	go func() {
		<-ctx.Done()
		s.logger.Warn("context canceled: shutting down")
		defer s.logger.Warn("shut down")
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed shutting down: %s", err)
		}
	}()

	s.logger.Info("listening on %s", s.address)
	err := server.ListenAndServe()
	if err != nil && ctx.Err() != context.Canceled { // server crashed
		crashed <- err
	}
}
