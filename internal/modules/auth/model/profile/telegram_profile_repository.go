package profile

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/pkg/repository"
)

type TelegramProfileRepository struct {
	*repository.BaseRepository
}

func NewTelegramProfileRepository(repository *repository.BaseRepository) *TelegramProfileRepository {
	return &TelegramProfileRepository{
		repository,
	}
}

type FindOptions struct {
	UserId int64
}

func (r TelegramProfileRepository) Find(options FindOptions, tx driver.Tx) (*TelegramProfile, error) {
	profile := TelegramProfile{}
	err := r.FindOneBy(
		goqu.Ex{"profile_id": goqu.Select("id").From("profiles").Where(goqu.Ex{"user_id": options.UserId})},
		&profile,
		tx,
	)
	if err != nil {
		return nil, fmt.Errorf("can not find telegram profile by options %v: %w", options, err)
	}
	if profile.Id == 0 {
		return nil, nil
	}

	return &profile, nil
}
