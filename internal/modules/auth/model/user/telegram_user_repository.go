package user

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
)

type TelegramUserRepository struct {
	*Repository
}

func NewTelegramUserRepository(repository *Repository) *TelegramUserRepository {
	return &TelegramUserRepository{repository}
}

type FindOptions struct {
	TelegramID int64
}

func (r TelegramUserRepository) Find(options FindOptions, tx driver.Tx) (*TelegramUser, error) {
	user := TelegramUser{}
	err := r.FindOneBy(goqu.Ex{"telegram_id": options.TelegramID}, &user, tx)
	if err != nil {
		return nil, fmt.Errorf("can not find telegram user by options %v: %w", options, err)
	}
	if user.Id == 0 {
		return nil, nil
	}

	return &user, nil
}
