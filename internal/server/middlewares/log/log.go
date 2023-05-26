package log

import (
	"net/http"
	"time"

	"github.com/qdm12/golibs/clientip"
)

func New(logger Logger, enabled bool) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &logHandler{
			childHandler:  handler,
			logger:        logger,
			enabled:       enabled,
			httpReqParser: clientip.NewParser(),
			timeNow:       time.Now,
		}
	}
}

type logHandler struct {
	childHandler  http.Handler
	logger        Logger
	enabled       bool
	httpReqParser clientip.HTTPRequestParser
	timeNow       func() time.Time
}

func (h *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.enabled {
		h.childHandler.ServeHTTP(w, r)
		return
	}
	startTime := h.timeNow()
	customWriter := &statefulWriter{ResponseWriter: w}
	h.childHandler.ServeHTTP(customWriter, r)
	clientIP := h.httpReqParser.ParseHTTPRequest(r)
	handlingDuration := h.timeNow().Sub(startTime).Round(time.Microsecond)
	h.logger.Infof("HTTP request: %d %s %s %s %dB %s",
		customWriter.status, r.Method, r.RequestURI, clientIP,
		customWriter.length, handlingDuration)
}
