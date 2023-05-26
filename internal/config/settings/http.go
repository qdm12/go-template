package settings

import (
	"fmt"
	"os"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
	"github.com/qdm12/govalid"
	"github.com/qdm12/govalid/address"
)

type HTTP struct {
	Address        *string
	RootURL        *string
	LogRequests    *bool
	AllowedOrigins []string
	AllowedHeaders []string
}

func (h *HTTP) setDefaults() {
	h.Address = gosettings.DefaultPointer(h.Address, ":8000")
	h.RootURL = gosettings.DefaultPointer(h.RootURL, "")
	h.LogRequests = gosettings.DefaultPointer(h.LogRequests, true)
	h.AllowedOrigins = gosettings.DefaultSliceRaw(h.AllowedOrigins, []string{})
	h.AllowedHeaders = gosettings.DefaultSliceRaw(h.AllowedHeaders, []string{})
}

func (h *HTTP) validate() (err error) {
	addressOption := address.OptionListening(os.Geteuid())
	err = govalid.ValidateAddress(*h.Address, addressOption)
	if err != nil {
		return fmt.Errorf("listening address: %w", err)
	}

	_, err = govalid.ValidateRootURL(*h.RootURL)
	if err != nil {
		fmt.Println("root url: ", *h.RootURL)
		return fmt.Errorf("root URL: %w", err)
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

	return node
}

func (h *HTTP) copy() (copied HTTP) {
	return HTTP{
		Address:        gosettings.CopyPointer(h.Address),
		RootURL:        gosettings.CopyPointer(h.RootURL),
		LogRequests:    gosettings.CopyPointer(h.LogRequests),
		AllowedOrigins: gosettings.CopySlice(h.AllowedOrigins),
		AllowedHeaders: gosettings.CopySlice(h.AllowedHeaders),
	}
}

func (h *HTTP) mergeWith(other HTTP) {
	h.Address = gosettings.MergeWithPointer(h.Address, other.Address)
	h.RootURL = gosettings.MergeWithPointer(h.RootURL, other.RootURL)
	h.LogRequests = gosettings.MergeWithPointer(h.LogRequests, other.LogRequests)
	h.AllowedOrigins = gosettings.MergeWithSlice(h.AllowedOrigins, other.AllowedOrigins)
	h.AllowedHeaders = gosettings.MergeWithSlice(h.AllowedHeaders, other.AllowedHeaders)
}

func (h *HTTP) overrideWith(other HTTP) {
	h.Address = gosettings.OverrideWithPointer(h.Address, other.Address)
	h.RootURL = gosettings.OverrideWithPointer(h.RootURL, other.RootURL)
	h.LogRequests = gosettings.OverrideWithPointer(h.LogRequests, other.LogRequests)
	h.AllowedOrigins = gosettings.OverrideWithSlice(h.AllowedOrigins, other.AllowedOrigins)
	h.AllowedHeaders = gosettings.OverrideWithSlice(h.AllowedHeaders, other.AllowedHeaders)
}
