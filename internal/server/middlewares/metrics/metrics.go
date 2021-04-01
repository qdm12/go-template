package metrics

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/qdm12/REPONAME_GITHUB/internal/metrics"
)

func New(metrics metrics.Metrics) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &metricsHandler{
			childHandler: handler,
			metrics:      metrics,
			timeNow:      time.Now,
		}
	}
}

type metricsHandler struct {
	childHandler http.Handler
	metrics      metrics.Metrics
	timeNow      func() time.Time // for mocks
}

func (h *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := h.timeNow()

	h.metrics.InflightRequestsGaugeAdd(1)
	defer h.metrics.InflightRequestsGaugeAdd(-1)

	statefulWriter := &statefulWriter{ResponseWriter: w}

	h.childHandler.ServeHTTP(statefulWriter, r)

	chiCtx := chi.RouteContext(r.Context())
	routePattern := chiCtx.RoutePattern()
	if routePattern == "" {
		routePattern = "unrecognized"
	}
	routePattern = strings.TrimSuffix(routePattern, "/")

	duration := h.timeNow().Sub(startTime)

	h.metrics.RequestCountInc(routePattern, statefulWriter.status)
	h.metrics.ResponseBytesCountAdd(routePattern, statefulWriter.status, statefulWriter.length)
	h.metrics.ResponseTimeHistogramObserve(routePattern, statefulWriter.status, duration)
}
