package contenttype

import (
	"fmt"
	"net/http"
	"strings"
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
		if responseContentType != "" {
			break
		}
	}

	if responseContentType == "" {
		responseContentType = JSON
		return "", responseContentType, fmt.Errorf("%w: %s", ErrRespContentTypeNotSupported, accept)
	}

	requestContentType = header.Get("Content-Type")
	requestContentType = strings.TrimSpace(requestContentType)
	if requestContentType == "" {
		requestContentType = JSON
	}
	if requestContentType != JSON {
		return "", responseContentType, fmt.Errorf("%w: %q", ErrContentTypeNotSupported, requestContentType)
	}

	return requestContentType, responseContentType, nil
}
