package log

import (
	"net/http"

	"github.com/qdm12/golibs/clientip"
	"github.com/qdm12/golibs/logging"
)

func New(logger logging.Logger, enabled bool) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &logHandler{
			childHandler: handler,
			logger:       logger,
			enabled:      enabled,
			ipExtractor:  clientip.NewExtractor(),
		}
	}
}

type logHandler struct {
	childHandler http.Handler
	logger       logging.Logger
	enabled      bool
	ipExtractor  clientip.Extractor
}

func (h *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.enabled {
		h.childHandler.ServeHTTP(w, r)
		return
	}
	customWriter := &statefulWriter{ResponseWriter: w}
	h.childHandler.ServeHTTP(customWriter, r)
	clientIP := h.ipExtractor.HTTPRequest(r)
	h.logger.Info("HTTP request: %d %s %s %s %d",
		customWriter.status, r.Method, r.RequestURI, clientIP, customWriter.length)
}
