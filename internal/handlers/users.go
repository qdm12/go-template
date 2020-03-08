package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/qdm12/REPONAME_GITHUB/internal/models"
	"github.com/qdm12/golibs/errors"
	"github.com/qdm12/golibs/server"
)

func (h *handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		ID uint64 `json:"id"`
	}
	if err := decoder.Decode(&body); err != nil {
		h.respondError(w, errors.NewBadRequest(err))
		return
	}
	user, err := h.proc.GetUserByID(body.ID)
	if err != nil {
		h.respondError(w, err)
		return
	}
	result := struct {
		User models.User `json:"user"`
	}{user}
	h.respondWrapper(w, server.Status(http.StatusOK), server.JSON(result))
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		User models.User `json:"user"`
	}
	if err := decoder.Decode(&body); err != nil {
		h.respondError(w, errors.NewBadRequest(err))
		return
	}
	if err := h.proc.CreateUser(body.User); err != nil {
		h.respondError(w, err)
		return
	}
	h.respondWrapper(w, server.Status(http.StatusOK))
}
