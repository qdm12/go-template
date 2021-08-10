package httperr

import (
	"net/http"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . ResponderInterface

type ResponderInterface interface {
	Respond(w http.ResponseWriter, status int, errString string)
}

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
