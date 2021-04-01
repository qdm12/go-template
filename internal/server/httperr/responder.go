package httperr

import (
	"net/http"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Reader

type Responder interface {
	Respond(w http.ResponseWriter, status int, errString string)
}

type responder struct {
	contentType string
}

func NewResponder(contentType string) Responder {
	return &responder{
		contentType: contentType,
	}
}

func (r *responder) Respond(w http.ResponseWriter, status int,
	errString string) {
	Respond(w, status, errString, r.contentType)
}
