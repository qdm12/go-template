package config

import (
	"github.com/qdm12/golibs/params"
)

type MemoryStore struct{}

func (m *MemoryStore) get(env params.Interface) (err error) {
	return nil
}
