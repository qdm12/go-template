package httperr

import (
	"net/http"
)

type Responder struct {
	contentType string
	logger      Logger
}

type Logger interface {
	Debugf(format string, args ...any)
}

func NewResponder(contentType string, logger Logger) *Responder {
	return &Responder{
		contentType: contentType,
		logger:      logger,
	}
}

// Respond responds the given error string and HTTP status
// to the given http response writer.
// If an error occurs responding, it is logged as a warning by
// the responder warner.
func (r *Responder) Respond(w http.ResponseWriter, status int,
	errString string,
) {
	err := Respond(w, status, errString, r.contentType)
	if err != nil {
		r.logger.Debugf("responding error: %s", err)
	}
}
