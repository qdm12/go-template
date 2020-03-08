package memory

import (
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
	"github.com/qdm12/golibs/errors"
)

func (db *Database) CreateUser(user models.User) (err error) {
	db.Lock()
	defer db.Unlock()
	db.data.Users = append(db.data.Users, user)
	return nil
}

func (db *Database) GetUserByID(id uint64) (user models.User, err error) {
	db.Lock()
	defer db.Unlock()
	for _, user := range db.data.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return user, errors.NewNotFound("user not found for id %d", id)
}
