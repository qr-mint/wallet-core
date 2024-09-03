package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	*repository.BaseRepository
}

func NewRepository(repository *repository.BaseRepository) *Repository {
	return &Repository{repository}
}

type FindOptions struct {
	UserId int64
}

func (r Repository) Find(options FindOptions, tx driver.Tx) (ProfileInterface, error) {
	profile := Profile{}
	err := r.FindOneBy(goqu.Ex{"user_id": options.UserId}, &profile, tx)
	if err != nil {
		return nil, fmt.Errorf("can not find profile by options %v: %w", options, err)
	}
	if profile.Id == 0 {
		return nil, nil
	}

	if profile.Type == app_enum.TelegramProfileType {
		telegramProfile := TelegramProfile{}
		err = r.FindOneBy(goqu.Ex{"profile_id": profile.Id}, &telegramProfile, tx)
		if err != nil {
			return nil, fmt.Errorf("can not find telegram profile by profile_id %v: %w", profile.Id, err)
		}
		if telegramProfile.Id == 0 {
			return nil, nil
		}
		telegramProfile.SetProfile(profile)

		return &telegramProfile, nil
	}

	return nil, fmt.Errorf("profile type is nit supported: %s", profile.Type)
}
