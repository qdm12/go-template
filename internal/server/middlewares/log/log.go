package log

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func New(logger Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &logHandler{
			childHandler: handler,
			logger:       logger,
			timeNow:      time.Now,
		}
	}
}

type logHandler struct {
	childHandler http.Handler
	logger       Logger
	timeNow      func() time.Time
}

func (h *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := h.timeNow()
	customWriter := &statefulWriter{ResponseWriter: w}
	h.childHandler.ServeHTTP(customWriter, r)
	clientIP := extractClientIP(r)
	bytesWritten := byteCountSI(customWriter.length)
	const durationResolution = time.Millisecond
	handlingDuration := h.timeNow().Sub(startTime).Round(durationResolution) / durationResolution
	h.logger.Infof("HTTP request: %d %s %s %s %s %dms",
		customWriter.status, r.Method, r.RequestURI, clientIP,
		bytesWritten, handlingDuration)
}

func byteCountSI(b int) string {
	const unit = 1000
	if b < unit {
		const base = 10
		return strconv.FormatInt(int64(b), base) + "B"
	}
	div := unit
	var exp uint
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
