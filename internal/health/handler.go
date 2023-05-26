package health

import (
	"net/http"
)

func NewHandler(logger Logger, healthcheck func() error) *Handler {
	return &Handler{
		logger:      logger,
		healthcheck: healthcheck,
	}
}

type Handler struct {
	logger      Logger
	healthcheck func() error
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
