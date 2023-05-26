package build

import (
	"encoding/json"
	"net/http"

	"github.com/qdm12/go-template/internal/server/contenttype"
	"github.com/qdm12/go-template/internal/server/httperr"
)

// Handler to get the program build information (GET /).
func (h *handler) getBuild(w http.ResponseWriter, r *http.Request) {
	_, responseContentType, err := contenttype.APICheck(r.Header)
	w.Header().Set("Content-Type", responseContentType)
	errResponder := httperr.NewResponder(responseContentType, h.logger)

	if err != nil {
		errResponder.Respond(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(h.build)
	if err != nil {
		h.logger.Error(err.Error())
		errResponder.Respond(w, http.StatusInternalServerError, "")
		return
	}
}
