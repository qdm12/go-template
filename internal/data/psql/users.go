package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/qdm12/REPONAME_GITHUB/internal/data/errors"
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
)

// CreateUser inserts a user in the database.
func (db *Database) CreateUser(ctx context.Context, user models.User) (err error) {
	_, err = db.sql.ExecContext(ctx,
		"INSERT INTO users(id, account, username, email) VALUES ($1,$2,$3,$4);",
		user.ID,
		user.Account,
		user.Username,
		user.Email,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrCreateUser, err)
	}
	return nil
}

// GetUserByID returns the user corresponding to a user ID from the database.
func (db *Database) GetUserByID(ctx context.Context, id uint64) (user models.User, err error) {
	row := db.sql.QueryRowContext(ctx,
		"SELECT account, email, username FROM users WHERE id = $1;",
		id,
	)
	user.ID = id
	err = row.Scan(&user.Account, &user.Email, &user.Username)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("%w: for id %d", errors.ErrUserNotFound, id)
	} else if err != nil {
		return user, fmt.Errorf("%w: for id %d: %s", errors.ErrGetUser, id, err)
	}
	return user, nil
}
