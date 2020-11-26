package psql

import (
	"database/sql"

	"github.com/qdm12/REPONAME_GITHUB/internal/models"
	"github.com/qdm12/golibs/errors"
)

// CreateUser inserts a user in the database.
func (db *Database) CreateUser(user models.User) (err error) {
	_, err = db.sql.Exec(
		"INSERT INTO users(id, account, username, email) VALUES ($1,$2,$3,$4);",
		user.ID,
		user.Account,
		user.Username,
		user.Email,
	)
	if err != nil {
		return errors.NewInternal("CreateUser: %s", err)
	}
	return nil
}

// GetUserByID returns the user corresponding to a user ID from the database.
func (db *Database) GetUserByID(id uint64) (user models.User, err error) {
	row := db.sql.QueryRow(
		"SELECT account, email, username FROM users WHERE id = $1;",
		id,
	)
	user.ID = id
	err = row.Scan(&user.Account, &user.Email, &user.Username)
	if err == sql.ErrNoRows {
		return user, errors.NewNotFound("no user found for id %d", id)
	} else if err != nil {
		return user, errors.NewInternal("GetUserByID for id %d: %s", id, err)
	}
	return user, nil
}
