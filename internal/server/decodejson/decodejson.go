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

// DecodeBody decodes the HTTP JSON encoded body into v.
// If the decoding succeeds, ok is returned as true.
// If the decoding fails, the function writes an error over
// HTTP to the client and returns ok as false.
// If writing the response to the client fails, an non-nil
// `errorResponseErr` error is returned as well.
// Therefore the caller should check for `errorResponseErr`
// every time `ok` is false.
func DecodeBody(w http.ResponseWriter, maxBytes int64,
	body io.ReadCloser, v any, responseContentType string) (
	ok bool, responseErr error,
) {
	body = http.MaxBytesReader(w, body, maxBytes)

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(v)
	if err != nil {
		errString, errCode := extractFromJSONErr(err)
		responseErr = httperr.Respond(w, errCode, errString, responseContentType)
		return false, responseErr
	}

	err = decoder.Decode(&struct{}{})
	if errors.Is(err, io.EOF) {
		return true, nil
	}
	const errString = "request body must only contain a single JSON object"
	responseErr = httperr.Respond(w, http.StatusBadRequest, errString, responseContentType)
	return false, responseErr
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
