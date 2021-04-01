package users

import (
	"net/http"

	"github.com/qdm12/go-template/internal/server/middlewares/cors"
)

func (h *handler) options(w http.ResponseWriter, r *http.Request) {
	cors.AllowCORSMethods(r, w, http.MethodPost, http.MethodGet)
}
