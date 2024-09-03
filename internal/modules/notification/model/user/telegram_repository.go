package user

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/pkg/repository"
)

type TelegramUserRepository struct {
	*repository.BaseRepository
}

func NewTelegramUserRepository(baseRepository *repository.BaseRepository) *TelegramUserRepository {
	return &TelegramUserRepository{baseRepository}
}

type FindOptions struct {
	UserId int64
}

func (r TelegramUserRepository) Find(options FindOptions, tx driver.Tx) (*TelegramUser, error) {
	user := TelegramUser{}
	err := r.FindOneBy(goqu.Ex{"user_id": options.UserId}, &user, tx)
	if err != nil {
		return nil, fmt.Errorf("can not find telegram user by options %v: %w", options, err)
	}
	if user.Id == 0 {
		return nil, nil
	}

	return &user, nil
}
