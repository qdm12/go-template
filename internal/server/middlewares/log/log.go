package log

import (
	"net/http"

	"github.com/qdm12/golibs/clientip"
)

func New(logger Logger, enabled bool) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &logHandler{
			childHandler:  handler,
			logger:        logger,
			enabled:       enabled,
			httpReqParser: clientip.NewParser(),
		}
	}
}

type logHandler struct {
	childHandler  http.Handler
	logger        Logger
	enabled       bool
	httpReqParser clientip.HTTPRequestParser
}

func (h *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.enabled {
		h.childHandler.ServeHTTP(w, r)
		return
	}
	customWriter := &statefulWriter{ResponseWriter: w}
	h.childHandler.ServeHTTP(customWriter, r)
	clientIP := h.httpReqParser.ParseHTTPRequest(r)
	h.logger.Infof("HTTP request: %d %s %s %s %dB",
		customWriter.status, r.Method, r.RequestURI, clientIP, customWriter.length)
}
