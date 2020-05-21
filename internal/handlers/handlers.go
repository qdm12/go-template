package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/qdm12/REPONAME_GITHUB/internal/processor"
	"github.com/qdm12/golibs/errors"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/server"
)

type handler struct {
	logger logging.Logger
	proc   processor.Processor
}

func NewHandler(rootURL string, proc processor.Processor, logger logging.Logger) http.HandlerFunc {
	h := &handler{
		proc:   proc,
		logger: logger,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// SOAP like api
		if r.Method != http.MethodPost {
			h.respondError(w, errors.NewBadRequest("HTTP method must be POST"))
			return
		}
		decoder := json.NewDecoder(r.Body)
		var body struct {
			Command string `json:"command"`
		}
		if err := decoder.Decode(&body); err != nil {
			h.respondError(w, errors.NewBadRequest(err))
			return
		}
		switch body.Command {
		case "get user by id":
			h.getUserByID(w, r)
		case "create user":
			h.createUser(w, r)
		default:
			h.respondError(w, errors.NewBadRequest("command %q is invalid", body.Command))
		}
	}
}

func (h *handler) respondWrapper(w http.ResponseWriter, setters ...server.ResponseSetter) {
	err := server.Respond(w, setters...)
	if err != nil {
		h.logger.Warn("cannot respond to client: %s", err)
	}
}

func (h *handler) respondError(w http.ResponseWriter, err error) {
	result := struct {
		Error string `json:"error"`
	}{"null"}
	if err != nil {
		result.Error = err.Error()
	}
	status := errors.HTTPStatus(err)
	h.respondWrapper(w, server.Status(status), server.JSON(result))
}
