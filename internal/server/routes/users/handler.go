// Package users is the HTTP handler for the users.
package users

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/qdm12/REPONAME_GITHUB/internal/processor"
	"github.com/qdm12/golibs/logging"
)

type handler struct {
	proc   processor.Processor
	logger logging.Logger
}

func NewHandler(logger logging.Logger, proc processor.Processor) http.Handler {
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
