package httperr

import (
	"net/http"
)

type Responder struct {
	contentType string
}

func NewResponder(contentType string) *Responder {
	return &Responder{
		contentType: contentType,
	}
}

func (r *Responder) Respond(w http.ResponseWriter, status int,
	errString string) {
	Respond(w, status, errString, r.contentType)
}
