package settings

import (
	"fmt"
	"os"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

type Websocket struct {
	Address *string
}

func (w *Websocket) setDefaults() {
	w.Address = gosettings.DefaultPointer(w.Address, ":8001")
}

func (w *Websocket) validate() (err error) {
	err = validate.ListeningAddress(*w.Address, os.Geteuid())
	if err != nil {
		return fmt.Errorf("listening address: %w", err)
	}

	return nil
}

func (w *Websocket) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Websocket server settings:")
	node.Appendf("Server listening address: %s", *w.Address)
	return node
}

func (w *Websocket) copy() (copied Websocket) {
	return Websocket{
		Address: gosettings.CopyPointer(w.Address),
	}
}

func (w *Websocket) overrideWith(other Websocket) {
	w.Address = gosettings.OverrideWithPointer(w.Address, other.Address)
}

func (w *Websocket) read(r *reader.Reader) {
	w.Address = r.Get("HTTP_WEBSOCKET_SERVER_ADDRESS")
}
