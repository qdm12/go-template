package users

import (
	"context"
	"errors"
	"net/http"

	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/go-template/internal/server/contenttype"
	"github.com/qdm12/go-template/internal/server/decodejson"
	"github.com/qdm12/go-template/internal/server/httperr"
)

// Handler for creating a user (POST /users/).
func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	_, responseContentType, err := contenttype.APICheck(r.Header)
	w.Header().Set("Content-Type", responseContentType)
	errResponder := httperr.NewResponder(responseContentType, h.logger)

	if err != nil {
		errResponder.Respond(w, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	ok, respondErr := decodejson.DecodeBody(w, 0, r.Body, &user, responseContentType)
	if !ok {
		if respondErr != nil {
			h.logger.Debugf("responding error: %s", respondErr)
		}
		return
	}

	if err := h.proc.CreateUser(r.Context(), user); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			errResponder.Respond(w, http.StatusRequestTimeout, "")
		} else {
			h.logger.Error(err.Error())
			errResponder.Respond(w, http.StatusInternalServerError, "")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	// TODO return ID created. ID to be set by data store, not by
	// client what the hell.
}
