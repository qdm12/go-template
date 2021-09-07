package config

import (
	"github.com/qdm12/golibs/params"
)

type Metrics struct {
	Address string
}

func (m *Metrics) get(env params.Interface) (warning string, err error) {
	m.Address, warning, err = m.getAddress(env)
	if err != nil {
		return warning, err
	}
	return warning, nil
}

func (m *Metrics) getAddress(env params.Interface) (address, warning string, err error) {
	const envKey = "METRICS_SERVER_ADDRESS"
	return env.ListeningAddress(envKey, params.Default(":9090"))
}
