package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/qdm12/REPONAME_GITHUB/internal/config"
	"github.com/qdm12/REPONAME_GITHUB/internal/metrics"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
	"github.com/qdm12/REPONAME_GITHUB/internal/processor"
	"github.com/qdm12/REPONAME_GITHUB/internal/server/middlewares/cors"
	logmware "github.com/qdm12/REPONAME_GITHUB/internal/server/middlewares/log"
	metricsmware "github.com/qdm12/REPONAME_GITHUB/internal/server/middlewares/metrics"
	"github.com/qdm12/REPONAME_GITHUB/internal/server/routes/build"
	"github.com/qdm12/REPONAME_GITHUB/internal/server/routes/users"
	"github.com/qdm12/golibs/logging"
)

func newRouter(config config.HTTP, logger logging.Logger,
	metrics metrics.Metrics, buildInfo models.BuildInformation,
	proc processor.Processor) http.Handler {
	router := chi.NewRouter()

	// Middlwares
	logMiddleware := logmware.New(logger, config.LogRequests)
	metricsMiddleware := metricsmware.New(metrics)
	corsMiddleware := cors.New(config.AllowedOrigins, config.AllowedHeaders)
	router.Use(metricsMiddleware, logMiddleware, corsMiddleware)

	APIPrefix := config.RootURL + "/api/v1"

	router.Mount(APIPrefix+"/users", users.NewHandler(logger, proc))
	router.Mount(APIPrefix+"/build", build.NewHandler(logger, buildInfo))

	// router.Handle("/metrics", promhttp.Handler()) // TODO

	return router
}
