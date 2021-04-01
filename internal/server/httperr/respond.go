package httperr

import (
	"encoding/json"
	"net/http"

	contenttype "github.com/qdm12/go-template/internal/server/contenttypes"
)

type errJSONWrapper struct {
	Error string `json:"error"`
}

func Respond(w http.ResponseWriter, status int,
	errString, contentType string) {
	w.WriteHeader(status)
	if errString == "" {
		errString = http.StatusText(status)
	}
	switch contentType {
	case contenttype.JSON:
		body := errJSONWrapper{Error: errString}
		_ = json.NewEncoder(w).Encode(body)
	default:
		_, _ = w.Write([]byte(errString))
	}
}
