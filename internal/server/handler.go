package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/qdm12/REPONAME_GITHUB/internal/models"
	"github.com/qdm12/REPONAME_GITHUB/internal/processor"
	"github.com/qdm12/golibs/logging"
)

func newHandler(rootURL string, logger logging.Logger,
	buildInfo models.BuildInformation, proc processor.Processor) http.Handler {
	return &handler{
		rootURL:   rootURL,
		logger:    logger,
		buildInfo: buildInfo,
		proc:      proc,
	}
}

type handler struct {
	rootURL   string
	logger    logging.Logger
	buildInfo models.BuildInformation
	proc      processor.Processor
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = strings.TrimPrefix(r.RequestURI, h.rootURL)

	// SOAP like API
	if !strings.HasPrefix(r.RequestURI, "/v1/") {
		httpError(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		httpError(w, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var body struct {
		Command string `json:"command"`
	}
	if err := decoder.Decode(&body); err != nil {
		httpBodyDecodeError(w, err)
		return
	}
	switch body.Command {
	case "get user by id":
		h.getUserByID(w, r)
	case "create user":
		h.createUser(w, r)
	default:
		http.Error(w, fmt.Sprintf("command %q is not valid", body.Command), http.StatusBadRequest)
	}
}
