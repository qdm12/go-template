package cors

import (
	"net/http"
)

func New(allowedOriginsSlice, allowedHeaders []string) func(handler http.Handler) http.Handler {
	allowedOrigins := make(map[string]struct{}, len(allowedOriginsSlice))
	for _, allowedOrigin := range allowedOriginsSlice {
		allowedOrigins[allowedOrigin] = struct{}{}
	}
	return func(handler http.Handler) http.Handler {
		return &corsHandler{
			childHandler:   handler,
			allowedOrigins: allowedOrigins,
			allowedHeaders: allowedHeaders,
		}
	}
}

type corsHandler struct {
	childHandler   http.Handler
	allowedOrigins map[string]struct{}
	allowedHeaders []string
}

func (h *corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = setCrossOriginHeaders(r.Header, w.Header(), h.allowedOrigins, h.allowedHeaders)
	// if error is not nil, CORS headers are NOT set so the browser will fail.
	// We don't want to stop the handling because it could be a request from a server.
	h.childHandler.ServeHTTP(w, r)
}
