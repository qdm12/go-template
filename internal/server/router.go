package server

import (
	"github.com/go-chi/chi"
	"github.com/qdm12/go-template/internal/config"
	"github.com/qdm12/go-template/internal/metrics"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/go-template/internal/processor"
	"github.com/qdm12/go-template/internal/server/middlewares/cors"
	logmware "github.com/qdm12/go-template/internal/server/middlewares/log"
	metricsmware "github.com/qdm12/go-template/internal/server/middlewares/metrics"
	"github.com/qdm12/go-template/internal/server/routes/build"
	"github.com/qdm12/go-template/internal/server/routes/users"
	"github.com/qdm12/golibs/logging"
)

func newRouter(config config.HTTP, logger logging.Logger,
	metrics metrics.Interface, buildInfo models.BuildInformation,
	proc processor.Interface) *chi.Mux {
	router := chi.NewRouter()

	// Middlewares
	logMiddleware := logmware.New(logger, config.LogRequests)
	metricsMiddleware := metricsmware.New(metrics)
	corsMiddleware := cors.New(config.AllowedOrigins, config.AllowedHeaders)
	router.Use(metricsMiddleware, logMiddleware, corsMiddleware)

	APIPrefix := config.RootURL + "/api/v1"

	router.Mount(APIPrefix+"/users", users.NewHandler(logger, proc))
	router.Mount(APIPrefix+"/build", build.NewHandler(logger, buildInfo))

	return router
}
