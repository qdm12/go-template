package processor

import (
	"context"

	"github.com/qdm12/REPONAME_GITHUB/internal/models"
)

func (p *processor) CreateUser(ctx context.Context, user models.User) error {
	return p.db.CreateUser(ctx, user)
}

func (p *processor) GetUserByID(ctx context.Context, id uint64) (user models.User, err error) {
	return p.db.GetUserByID(ctx, id)
}
