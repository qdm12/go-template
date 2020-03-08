package processor

import (
	"github.com/qdm12/REPONAME_GITHUB/internal/models"
)

func (p *processor) CreateUser(user models.User) error {
	return p.db.CreateUser(user)
}

func (p *processor) GetUserByID(id uint64) (user models.User, err error) {
	user, err = p.db.GetUserByID(id)
	if err != nil {
		return user, err
	}
	return user, nil
}
