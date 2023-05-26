package server

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/qdm12/go-template/internal/config/settings"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/go-template/internal/server/middlewares/cors"
	logmware "github.com/qdm12/go-template/internal/server/middlewares/log"
	metricsmware "github.com/qdm12/go-template/internal/server/middlewares/metrics"
	"github.com/qdm12/go-template/internal/server/routes/build"
	"github.com/qdm12/go-template/internal/server/routes/users"
)

func NewRouter(config settings.HTTP, logger Logger,
	metrics Metrics, buildInfo models.BuildInformation,
	proc Processor) *chi.Mux {
	router := chi.NewRouter()

	var middlewares []func(http.Handler) http.Handler
	metricsMiddleware := metricsmware.New(metrics)
	middlewares = append(middlewares, metricsMiddleware)
	if *config.LogRequests {
		logMiddleware := logmware.New(logger)
		middlewares = append(middlewares, logMiddleware)
	}
	corsMiddleware := cors.New(config.AllowedOrigins, config.AllowedHeaders)
	middlewares = append(middlewares, corsMiddleware)
	router.Use(middlewares...)

	APIPrefix := *config.RootURL
	for strings.HasSuffix(APIPrefix, "/") {
		APIPrefix = strings.TrimSuffix(APIPrefix, "/")
	}
	APIPrefix += "/api/v1"

	router.Mount(APIPrefix+"/users", users.NewHandler(logger, proc))
	router.Mount(APIPrefix+"/build", build.NewHandler(logger, buildInfo))

	return router
}
