package json

import (
	"context"
	"fmt"

	"github.com/qdm12/go-template/internal/data/errors"
	"github.com/qdm12/go-template/internal/models"
)

func (db *Database) CreateUser(ctx context.Context, user models.User) (err error) {
	if err := db.memory.CreateUser(ctx, user); err != nil {
		return err
	}
	if err := db.writeFile(); err != nil {
		return fmt.Errorf("%w: for user %#v: %s", errors.ErrCreateUser, user, err)
	}
	return nil
}

func (db *Database) GetUserByID(ctx context.Context, id uint64) (user models.User, err error) {
	return db.memory.GetUserByID(ctx, id)
}
