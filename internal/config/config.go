// Package config takes care of reading and checking the program configuration
// from environment variables.
package config

import (
	"errors"
	"fmt"

	"github.com/qdm12/golibs/params"
)

type Config struct {
	HTTP    HTTP
	Metrics Metrics
	Log     Log
	Store   Store
	Health  Health
}

var (
	ErrHTTPConfig    = errors.New("cannot obtain HTTP server config")
	ErrStoreConfig   = errors.New("cannot obtain store config")
	ErrHealthConfig  = errors.New("cannot obtain health config")
	ErrMetricsConfig = errors.New("cannot obtain metrics config")
	ErrLogConfig     = errors.New("cannot obtain log config")
)

func (c *Config) Read(env params.Interface) (warnings []string, err error) {
	warning, err := c.HTTP.get(env)
	if len(warning) > 0 {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrHTTPConfig, err)
	}

	warning, err = c.Store.get(env)
	if len(warning) > 0 {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrStoreConfig, err)
	}

	warning, err = c.Health.get(env)
	if len(warning) > 0 {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrHealthConfig, err)
	}

	warning, err = c.Metrics.get(env)
	if len(warning) > 0 {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrMetricsConfig, err)
	}

	err = c.Log.get(env)
	if err != nil {
		return warnings, fmt.Errorf("%w: %s", ErrLogConfig, err)
	}

	return warnings, nil
}
