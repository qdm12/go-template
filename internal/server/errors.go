package server

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBodyDecode = errors.New("cannot decode request body")
)

func httpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func httpBodyDecodeError(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("%s: %s", ErrBodyDecode, err), http.StatusBadRequest)
}
