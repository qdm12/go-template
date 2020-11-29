package json

import (
	"fmt"

	"github.com/qdm12/REPONAME_GITHUB/internal/data/errors"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
)

func (db *Database) CreateUser(user models.User) (err error) {
	if err := db.memory.CreateUser(user); err != nil {
		return err
	}
	if err := db.writeFile(); err != nil {
		return fmt.Errorf("%w: for user %#v: %s", errors.ErrCreateUser, user, err)
	}
	return nil
}

func (db *Database) GetUserByID(id uint64) (user models.User, err error) {
	return db.memory.GetUserByID(id)
}
