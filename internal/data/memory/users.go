package memory

import (
	"context"
	"fmt"

	"github.com/qdm12/go-template/internal/data/errors"
	"github.com/qdm12/go-template/internal/models"
)

func (db *Database) CreateUser(_ context.Context, user models.User) (err error) {
	db.Lock()
	defer db.Unlock()
	db.data.Users = append(db.data.Users, user)
	return nil
}

func (db *Database) GetUserByID(_ context.Context, id uint64) (user models.User, err error) {
	db.Lock()
	defer db.Unlock()
	for _, user := range db.data.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return user, fmt.Errorf("%w: for id %d", errors.ErrUserNotFound, id)
}
