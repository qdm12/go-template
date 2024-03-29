package settings

import (
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

type MemoryDatabase struct{}

func (m *MemoryDatabase) setDefaults() {}

func (m *MemoryDatabase) validate() (err error) { return nil }

func (m *MemoryDatabase) toLinesNode() (node *gotree.Node) {
	return nil
}

func (m *MemoryDatabase) copy() (copied MemoryDatabase) {
	return MemoryDatabase{}
}

func (m *MemoryDatabase) overrideWith(MemoryDatabase) {}

func (m *MemoryDatabase) read(_ *reader.Reader) {}
