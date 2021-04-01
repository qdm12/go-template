package server

import (
	"encoding/json"
	"errors"
	"net/http"

	dataerr "github.com/qdm12/REPONAME_GITHUB/internal/data/errors"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
)

func (h *handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		ID uint64 `json:"id"`
	}
	if err := decoder.Decode(&body); err != nil {
		httpBodyDecodeError(w, err)
		return
	}
	user, err := h.proc.GetUserByID(r.Context(), body.ID)
	if err != nil {
		switch {
		case errors.Is(err, dataerr.ErrGetUser):
			h.logger.Error(err)
			httpError(w, http.StatusInternalServerError)
		case errors.Is(err, dataerr.ErrUserNotFound):
			httpError(w, http.StatusNotFound)
		default:
			h.logger.Error(err)
			httpError(w, http.StatusInternalServerError)
		}
		return
	}
	result := struct {
		User models.User `json:"user"`
	}{user}
	b, err := json.Marshal(result)
	if err != nil {
		h.logger.Error(err)
		httpError(w, http.StatusInternalServerError)
	}
	if _, err := w.Write(b); err != nil {
		h.logger.Error(err)
	}
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		User models.User `json:"user"`
	}
	if err := decoder.Decode(&body); err != nil {
		httpBodyDecodeError(w, err)
		return
	}
	if err := h.proc.CreateUser(r.Context(), body.User); err != nil {
		switch {
		case errors.Is(err, dataerr.ErrCreateUser):
			h.logger.Error(err)
			httpError(w, http.StatusInternalServerError)
		default:
			h.logger.Error(err)
			httpError(w, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
