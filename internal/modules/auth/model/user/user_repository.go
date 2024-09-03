package user

import (
	"database/sql/driver"
	"fmt"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	*repository.BaseRepository
}

func NewRepository(repository *repository.BaseRepository) *Repository {
	return &Repository{repository}
}

func (r Repository) Save(user UserInterface, tx driver.Tx) error {
	baseUser := user.GetUser()
	err := r.CreateOrUpdate(&baseUser, tx)
	if err != nil {
		return fmt.Errorf("can not create user: %s", err)
	}

	user.SetUser(baseUser)
	err = r.CreateOrUpdate(user, tx)
	if err != nil {
		return fmt.Errorf("can not create user: %s", err)
	}
	err = r.Refresh(user, tx)
	if err != nil {
		return fmt.Errorf("can not refresh user: %s", err)
	}

	return nil
}
