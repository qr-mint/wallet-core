package profile

import (
	"database/sql/driver"
	"fmt"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	*repository.BaseRepository
}

func NewRepository(repository *repository.BaseRepository) *Repository {
	return &Repository{
		repository,
	}
}

func (r Repository) Save(profile ProfileInterface, tx driver.Tx) error {
	if profile.GetId() == 0 {
		baseProfile := profile.GetProfile()
		err := r.CreateOrUpdate(&baseProfile, tx)
		if err != nil {
			return fmt.Errorf("can not save profile: %s", err)
		}
		profile.SetProfile(baseProfile)
	}

	err := r.CreateOrUpdate(profile, tx)
	if err != nil {
		return fmt.Errorf("can not save profile: %s", err)
	}
	err = r.Refresh(profile, tx)
	if err != nil {
		return fmt.Errorf("can not refresh profile: %s", err)
	}

	return nil
}
