package health

import (
	"net/http"
)

func NewHandler(logger Logger, healthcheck func() error) http.Handler {
	return &handler{
		logger:      logger,
		healthcheck: healthcheck,
	}
}

type handler struct {
	logger      Logger
	healthcheck func() error
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || (r.RequestURI != "" && r.RequestURI != "/") {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err := h.healthcheck(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
