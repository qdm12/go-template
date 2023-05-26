package env

import (
	"github.com/qdm12/go-template/internal/config/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func readMetrics() (metrics settings.Metrics) {
	metrics.Address = env.Get("METRICS_SERVER_ADDRESS")
	return metrics
}
