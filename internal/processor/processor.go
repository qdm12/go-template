package processor

import (
	"github.com/qdm12/REPONAME_GITHUB/internal/data"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"

	"github.com/qdm12/golibs/crypto"
)

// Processor has methods to process data and return results
type Processor interface {
	GetUserByID(id uint64) (user models.User, err error)
	CreateUser(user models.User) (err error)
}

type processor struct {
	db     data.Database
	crypto crypto.Crypto
}

// NewProcessor creates a new processor object
func NewProcessor(db data.Database, crypto crypto.Crypto) Processor {
	return &processor{db, crypto}
}
