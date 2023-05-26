package env

import (
	"github.com/qdm12/go-template/internal/config/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) ReadHealth() (health settings.Health) {
	health.Address = env.Get("HEALTH_ADDRESS")
	return health
}
