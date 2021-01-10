package config

import (
	"github.com/qdm12/golibs/params"
)

type HTTP struct {
	Address string
	RootURL string
}

func (h *HTTP) get(env params.Env) (warning string, err error) {
	h.Address, warning, err = h.getAddress(env)
	if err != nil {
		return warning, err
	}
	h.RootURL, err = h.getRootURL(env)
	if err != nil {
		return warning, err
	}
	return warning, nil
}

func (h *HTTP) getAddress(env params.Env) (address, warning string, err error) {
	const envKey = "HTTP_SERVER_ADDRESS"
	options := []params.OptionSetter{
		params.Default("0.0.0.0:8000"),
	}
	return env.ListeningAddress(envKey, options...)
}

func (h *HTTP) getRootURL(env params.Env) (rootURL string, err error) {
	return env.RootURL("HTTP_SERVER_ROOT_URL")
}
