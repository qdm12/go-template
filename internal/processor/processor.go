package processor

import (
	"context"

	"github.com/qdm12/REPONAME_GITHUB/internal/data"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
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
