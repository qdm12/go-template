// Package processor contains operations the server can run and
// serves as the middle ground between the network server and
// the data store.
package processor

import (
	"context"

	"github.com/qdm12/go-template/internal/data"
	"github.com/qdm12/go-template/internal/models"
	"github.com/qdm12/golibs/crypto"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Processor

// Processor has methods to process data and return results.
type Processor interface {
	GetUserByID(ctx context.Context, id uint64) (user models.User, err error)
	CreateUser(ctx context.Context, user models.User) (err error)
}

type processor struct {
	db     data.Database
	crypto crypto.Crypto
}

// NewProcessor creates a new processor object.
func NewProcessor(db data.Database, crypto crypto.Crypto) Processor {
	return &processor{db, crypto}
}
