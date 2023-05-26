package processor

import (
	"context"
	"errors"
	"fmt"

	dataerr "github.com/qdm12/go-template/internal/data/errors"
	"github.com/qdm12/go-template/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (p *Processor) CreateUser(ctx context.Context, user models.User) error {
	return p.db.CreateUser(ctx, user)
}

func (p *Processor) GetUserByID(ctx context.Context, id uint64) (user models.User, err error) {
	user, err = p.db.GetUserByID(ctx, id)
	if errors.Is(err, dataerr.ErrUserNotFound) {
		err = fmt.Errorf("%w: %s", ErrUserNotFound, errors.Unwrap(err))
	}
	return user, err
}
