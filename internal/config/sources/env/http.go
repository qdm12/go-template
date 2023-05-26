package env

import (
	"fmt"

	"github.com/qdm12/go-template/internal/config/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readHTTP() (http settings.HTTP, err error) {
	http.Address = env.StringPtr("HTTP_SERVER_ADDRESS")
	http.RootURL = env.StringPtr("HTTP_SERVER_ROOT_URL")
	http.LogRequests, err = env.BoolPtr("HTTP_SERVER_LOG_REQUESTS")
	if err != nil {
		return http, fmt.Errorf("environment variable HTTP_SERVER_LOG_REQUESTS: %w", err)
	}
	http.AllowedOrigins = env.CSV("HTTP_SERVER_ALLOWED_ORIGINS")
	http.AllowedHeaders = env.CSV("HTTP_SERVER_ALLOWED_HEADERS")
	return http, nil
}
