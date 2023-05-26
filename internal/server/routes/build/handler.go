// Package build is the HTTP handler for the build information.
package build

import (
	"github.com/go-chi/chi"
	"github.com/qdm12/go-template/internal/models"
)

type handler struct {
	logger Logger
	build  models.BuildInformation
}

func NewHandler(logger Logger, buildInfo models.BuildInformation) *chi.Mux {
	h := &handler{
		logger: logger,
		build:  buildInfo,
	}
	router := chi.NewRouter()
	router.Get("/", h.getBuild)
	router.Options("/", h.options)
	return router
}
