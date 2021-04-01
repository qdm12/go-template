// Package cors has a middleware and functions to parse and set
// CORS correctly.
package cors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	errNoOriginHeader   = errors.New("no origin header found")
	errOriginNotAllowed = errors.New("origin not allowed")
)

func setCrossOriginHeaders(requestHeaders, responseHeaders http.Header,
	allowedOrigins map[string]struct{}, allowedHeaders []string) error {
	origin := requestHeaders.Get("Origin")
	if len(origin) == 0 {
		return errNoOriginHeader
	}

	if _, ok := allowedOrigins[origin]; !ok {
		return fmt.Errorf("%w: %s", errOriginNotAllowed, origin)
	}

	responseHeaders.Set("Access-Control-Allow-Origin", origin)
	responseHeaders.Set("Access-Control-Max-Age", "14400") // 4 hours
	for i := range allowedHeaders {
		responseHeaders.Add("Access-Control-Allow-Headers", allowedHeaders[i])
	}

	return nil
}

func AllowCORSMethods(r *http.Request, w http.ResponseWriter, methods ...string) {
	methods = append(methods, http.MethodOptions)
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))

	requestMethod := r.Header.Get("Access-Control-Request-Method")
	for _, method := range methods {
		if method == requestMethod {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
