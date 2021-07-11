package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	dataerr "github.com/qdm12/go-template/internal/data/errors"
	contenttype "github.com/qdm12/go-template/internal/server/contenttypes"
	"github.com/qdm12/go-template/internal/server/httperr"
)

// Handler to get a user by ID (GET /users/{id}).
func (h *handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	_, responseContentType, err := contenttype.APICheck(r.Header)
	w.Header().Set("Content-Type", responseContentType)
	errResponder := httperr.NewResponder(responseContentType)

	if err != nil {
		errResponder.Respond(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := extractUserID(r)
	if err != nil {
		errResponder.Respond(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.proc.GetUserByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, dataerr.ErrUserNotFound):
			errResponder.Respond(w, http.StatusNotFound, err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			errResponder.Respond(w, http.StatusRequestTimeout, "")
		default:
			h.logger.Error(err)
			errResponder.Respond(w, http.StatusInternalServerError, "")
		}
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		h.logger.Error(err)
		errResponder.Respond(w, http.StatusInternalServerError, "")
		return
	}
}

var (
	errUserIDMissingURLPath = errors.New("user ID must be provided in the URL path")
	errUserIDMalformed      = errors.New("user ID is malformed")
)

func extractUserID(r *http.Request) (id uint64, err error) {
	s := chi.URLParam(r, "id")
	if s == "" {
		return 0, errUserIDMissingURLPath
	}
	id, err = strconv.ParseUint(s, 10, 64) //nolint:gomnd
	if err != nil {
		return 0, fmt.Errorf("%w: %q", errUserIDMalformed, s)
	}
	return id, nil
}
