package json

import (
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
	"github.com/qdm12/golibs/errors"
)

func (db *Database) CreateUser(user models.User) (err error) {
	if err := db.memory.CreateUser(user); err != nil {
		return err
	}
	if err := db.writeFile(); err != nil {
		return errors.NewInternal("CreateUser for %#v: %s", user, err)
	}
	return nil
}

func (db *Database) GetUserByID(id uint64) (user models.User, err error) {
	return db.memory.GetUserByID(id)
}
