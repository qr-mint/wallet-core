package user

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type User struct {
	Id   int64             `primary:"true" must_generate:"true" db:"id"`
	Type app_enum.UserType `db:"type"`
}

func (u User) GetTableName() string {
	return "users"
}

func (u User) Clear() repository.Model {
	return &User{}
}

type TelegramUser struct {
	Id         int64 `primary:"true" must_generate:"true" db:"id"`
	TelegramId int64 `db:"telegram_id"`
}

func (u *TelegramUser) GetTableName() string {
	return "telegram_users"
}
