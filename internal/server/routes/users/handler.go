// Package users is the HTTP handler for the users.
package users

import (
	"github.com/go-chi/chi"
)

type handler struct {
	proc   Processor
	logger Logger
}

func NewHandler(logger Logger, proc Processor) *chi.Mux {
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
