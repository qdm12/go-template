// Package processor contains operations the server can run and
// serves as the middle ground between the network server and
// the data store.
package processor

import (
	"github.com/qdm12/go-template/internal/data"
)

type Processor struct {
	db data.Database
}

// NewProcessor creates a new Processor object.
func NewProcessor(db data.Database) *Processor {
	return &Processor{
		db: db,
	}
}
