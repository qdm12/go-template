package settings

import (
	"fmt"
	"os"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

type Metrics struct {
	Address string
}

func (m *Metrics) setDefaults() {
	m.Address = gosettings.DefaultComparable(m.Address, ":9090")
}

func (m *Metrics) validate() (err error) {
	err = validate.ListeningAddress(m.Address, os.Geteuid())
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

func (m *Metrics) overrideWith(other Metrics) {
	m.Address = gosettings.OverrideWithComparable(m.Address, other.Address)
}

func (m *Metrics) read(r *reader.Reader) {
	m.Address = r.String("METRICS_SERVER_ADDRESS")
}
