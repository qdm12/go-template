package settings

import (
	"fmt"
	"os"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
	"github.com/qdm12/govalid"
	"github.com/qdm12/govalid/address"
)

type Health struct {
	Address string
}

func (h *Health) SetDefaults() {
	h.Address = "127.0.0.1:9999"
}

func (h *Health) Validate() (err error) {
	addressOption := address.OptionListening(os.Geteuid())
	err = govalid.ValidateAddress(h.Address, addressOption)
	if err != nil {
		return fmt.Errorf("listening address: %w", err)
	}
	return nil
}

func (h *Health) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Health settings:")
	node.Appendf("Server listening address: %s", h.Address)
	return node
}

func (h *Health) copy() (copied Health) {
	return Health{
		Address: h.Address,
	}
}

func (h *Health) mergeWith(other Health) {
	h.Address = gosettings.MergeWithString(h.Address, other.Address)
}

func (h *Health) overrideWith(other Health) {
	h.Address = gosettings.OverrideWithString(h.Address, other.Address)
}
