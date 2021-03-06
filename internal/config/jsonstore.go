package config

import (
	"github.com/qdm12/golibs/params"
)

type JSONStore struct {
	Filepath string
}

func (j *JSONStore) get(env params.Env) (err error) {
	return j.getFilepath(env)
}

func (j *JSONStore) getFilepath(env params.Env) (err error) {
	const envKey = "STORE_JSON_FILEPATH"
	j.Filepath, err = env.Path(envKey, params.Default("data.json"))
	return err
}
