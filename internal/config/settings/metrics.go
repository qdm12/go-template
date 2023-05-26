package settings

import (
	"fmt"
	"os"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
	"github.com/qdm12/govalid"
	"github.com/qdm12/govalid/address"
)

type Metrics struct {
	Address string
}

func (m *Metrics) setDefaults() {
	m.Address = gosettings.DefaultString(m.Address, ":9090")
}

func (m *Metrics) validate() (err error) {
	addressOption := address.OptionListening(os.Geteuid())
	err = govalid.ValidateAddress(m.Address, addressOption)
	if err != nil {
		return fmt.Errorf("listening address: %w", err)
	}

	return nil
}

func (m *Metrics) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Metrics settings:")
	node.Appendf("Server listening address: %s", m.Address)
	return node
}

func (m *Metrics) copy() (copied Metrics) {
	return Metrics{
		Address: m.Address,
	}
}

func (m *Metrics) mergeWith(other Metrics) {
	m.Address = gosettings.MergeWithString(m.Address, other.Address)
}

func (m *Metrics) overrideWith(other Metrics) {
	m.Address = gosettings.OverrideWithString(m.Address, other.Address)
}
