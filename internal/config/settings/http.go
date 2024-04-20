package settings

import (
	"fmt"
	"os"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

type HTTP struct {
	Address        *string
	RootURL        *string
	LogRequests    *bool
	AllowedOrigins []string
	AllowedHeaders []string
	Websocket      Websocket
}

func (h *HTTP) setDefaults() {
	h.Address = gosettings.DefaultPointer(h.Address, ":8000")
	h.RootURL = gosettings.DefaultPointer(h.RootURL, "")
	h.LogRequests = gosettings.DefaultPointer(h.LogRequests, true)
	h.AllowedOrigins = gosettings.DefaultSlice(h.AllowedOrigins, []string{})
	h.AllowedHeaders = gosettings.DefaultSlice(h.AllowedHeaders, []string{})
	h.Websocket.setDefaults()
}

func (h *HTTP) validate() (err error) {
	err = validate.ListeningAddress(*h.Address, os.Geteuid())
	if err != nil {
		return fmt.Errorf("listening address: %w", err)
	}

	err = h.Websocket.validate()
	if err != nil {
		return fmt.Errorf("websocket: %w", err)
	}

	return nil
}

func (h *HTTP) toLinesNode() (node *gotree.Node) {
	node = gotree.New("HTTP server settings:")
	node.Appendf("Server listening address: %s", *h.Address)
	node.Appendf("Root URL: %s", *h.RootURL)
	node.Appendf("Log requests: %s", boolPtrToYesNo(h.LogRequests))

	allowedOriginsNode := gotree.New("Allowed origins:")
	for _, allowedOrigin := range h.AllowedOrigins {
		allowedOriginsNode.Appendf(allowedOrigin)
	}
	node.AppendNode(allowedOriginsNode)

	allowedHeadersNode := gotree.New("Allowed headers:")
	for _, allowedHeader := range h.AllowedHeaders {
		allowedHeadersNode.Appendf(allowedHeader)
	}
	node.AppendNode(allowedHeadersNode)

	node.AppendNode(h.Websocket.toLinesNode())

	return node
}

func (h *HTTP) copy() (copied HTTP) {
	return HTTP{
		Address:        gosettings.CopyPointer(h.Address),
		RootURL:        gosettings.CopyPointer(h.RootURL),
		LogRequests:    gosettings.CopyPointer(h.LogRequests),
		AllowedOrigins: gosettings.CopySlice(h.AllowedOrigins),
		AllowedHeaders: gosettings.CopySlice(h.AllowedHeaders),
		Websocket:      h.Websocket.copy(),
	}
}

func (h *HTTP) overrideWith(other HTTP) {
	h.Address = gosettings.OverrideWithPointer(h.Address, other.Address)
	h.RootURL = gosettings.OverrideWithPointer(h.RootURL, other.RootURL)
	h.LogRequests = gosettings.OverrideWithPointer(h.LogRequests, other.LogRequests)
	h.AllowedOrigins = gosettings.OverrideWithSlice(h.AllowedOrigins, other.AllowedOrigins)
	h.AllowedHeaders = gosettings.OverrideWithSlice(h.AllowedHeaders, other.AllowedHeaders)
	h.Websocket.overrideWith(other.Websocket)
}

func (h *HTTP) read(r *reader.Reader) (err error) {
	h.Address = r.Get("HTTP_SERVER_ADDRESS")
	h.RootURL = r.Get("HTTP_SERVER_ROOT_URL")
	h.LogRequests, err = r.BoolPtr("HTTP_SERVER_LOG_REQUESTS")
	if err != nil {
		return fmt.Errorf("environment variable HTTP_SERVER_LOG_REQUESTS: %w", err)
	}
	h.AllowedOrigins = r.CSV("HTTP_SERVER_ALLOWED_ORIGINS")
	h.AllowedHeaders = r.CSV("HTTP_SERVER_ALLOWED_HEADERS")
	h.Websocket.read(r)
	return nil
}
