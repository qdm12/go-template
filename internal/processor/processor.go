// Package processor contains operations the server can run and
// serves as the middle ground between the network server and
// the data store.
package processor

type Processor struct {
	db Database
}

// NewProcessor creates a new Processor object.
func NewProcessor(db Database) *Processor {
	return &Processor{
		db: db,
	}
}
