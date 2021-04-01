// Package decodejson has helper functions to decode HTTP bodies encoded in JSON.
package decodejson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/qdm12/go-template/internal/server/httperr"
)

func DecodeBody(w http.ResponseWriter, maxBytes int64,
	body io.ReadCloser, v interface{}, responseContentType string) (ok bool) {
	body = http.MaxBytesReader(w, body, maxBytes)

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(v)
	if err != nil {
		errString, errCode := extractFromJSONErr(err)
		httperr.Respond(w, errCode, errString, responseContentType)
		return false
	}
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		const errString = "request body must only contain a single JSON object"
		httperr.Respond(w, http.StatusBadRequest, errString, responseContentType)
		return false
	}
	return true
}

func extractFromJSONErr(err error) (errString string, errCode int) {
	var (
		syntaxError        *json.SyntaxError
		unmarshalTypeError *json.UnmarshalTypeError
	)
	switch {
	case errors.As(err, &syntaxError):
		const format = "request body contains badly-formed JSON (at position %d)"
		return fmt.Sprintf(format, syntaxError.Offset), http.StatusBadRequest

	case errors.Is(err, io.ErrUnexpectedEOF):
		return "request body contains badly-formed JSON", http.StatusBadRequest

	case errors.As(err, &unmarshalTypeError):
		const format = "request body contains an invalid value for the %q field (at position %d)"
		return fmt.Sprintf(format, unmarshalTypeError.Field, unmarshalTypeError.Offset), http.StatusBadRequest

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		const format = "request body contains unknown field %s"
		return fmt.Sprintf(format, fieldName), http.StatusBadRequest

	case errors.Is(err, io.EOF):
		return "request body cannot be empty", http.StatusBadRequest

	case err.Error() == "http: request body too large":
		return "request body is too large", http.StatusRequestEntityTooLarge

	default:
		return err.Error(), http.StatusInternalServerError
	}
}
