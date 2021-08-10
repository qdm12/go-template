// Package users is the HTTP handler for the users.
package users

import (
	"github.com/go-chi/chi"
	"github.com/qdm12/go-template/internal/processor"
	"github.com/qdm12/golibs/logging"
)

type handler struct {
	proc   processor.Interface
	logger logging.Logger
}

func NewHandler(logger logging.Logger, proc processor.Interface) *chi.Mux {
	h := &handler{
		proc:   proc,
		logger: logger,
	}
	router := chi.NewRouter()
	router.Get("/{id}", h.getUserByID)
	router.Post("/", h.createUser)
	router.Options("/", h.options)
	return router
}
