package httperr

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qdm12/go-template/internal/server/contenttype"
)

type errJSONWrapper struct {
	Error string `json:"error"`
}

func Respond(w http.ResponseWriter, status int,
	errString, contentType string,
) (err error) {
	w.WriteHeader(status)
	if errString == "" {
		errString = http.StatusText(status)
	}
	switch contentType {
	case contenttype.JSON:
		body := errJSONWrapper{Error: errString}
		err = json.NewEncoder(w).Encode(body)
		if err != nil {
			return fmt.Errorf("encoding and writing JSON response: %w", err)
		}
	default:
		_, err = w.Write([]byte(errString))
		if err != nil {
			return fmt.Errorf("writing raw error string: %w", err)
		}
	}
	return nil
}
