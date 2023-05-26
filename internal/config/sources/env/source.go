package env

import (
	"fmt"

	"github.com/qdm12/go-template/internal/config/settings"
)

type Source struct{}

func New() *Source {
	return &Source{}
}

func (s *Source) String() string { return "environment variables" }

func (s *Source) Read() (settings settings.Settings, err error) {
	settings.HTTP, err = s.readHTTP()
	if err != nil {
		return settings, fmt.Errorf("HTTP server settings: %w", err)
	}

	settings.Metrics = readMetrics()
	settings.Log, err = readLog()
	if err != nil {
		return settings, fmt.Errorf("logging settings: %w", err)
	}
	settings.Database = readDatabase()
	settings.Health = s.ReadHealth()

	return settings, nil
}
