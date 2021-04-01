package contenttype

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	JSON = "application/json"
	HTML = "text/html"
)

var (
	errContentTypeNotSupported     = errors.New("content type is not supported")
	errRespContentTypeNotSupported = errors.New("no response content type supported")
)

func APICheck(header http.Header) (
	requestContentType, responseContentType string, err error) {
	accept := header.Get("Accept")
	if accept == "" {
		accept = JSON
	}

	acceptedTypes := strings.Split(accept, ",")
	for _, acceptedType := range acceptedTypes {
		acceptedType = strings.TrimSpace(acceptedType)
		switch acceptedType {
		case JSON:
			responseContentType = acceptedType
		case HTML:
			responseContentType = JSON // override for browser access to the API
		}
		if len(responseContentType) > 0 {
			break
		}
	}

	if len(responseContentType) == 0 {
		responseContentType = JSON
		return "", responseContentType, fmt.Errorf("%w: %s", errRespContentTypeNotSupported, accept)
	}

	requestContentType = header.Get("Content-Type")
	requestContentType = strings.TrimSpace(requestContentType)
	if requestContentType == "" {
		requestContentType = JSON
	}
	if requestContentType != JSON {
		return "", responseContentType, fmt.Errorf("%w: %q", errContentTypeNotSupported, requestContentType)
	}

	return requestContentType, responseContentType, nil
}
