package user

import (
	"nexus-wallet/internal/app_enum"
)

type UserInterface interface {
	GetUser() User
	SetUser(User)
	GetTableName() string
}

type User struct {
	Id   int64             `primary:"true" must_generate:"true" db:"id"`
	Type app_enum.UserType `db:"type"`
}

func (u User) GetTableName() string {
	return "users"
}

type TelegramUser struct {
	Id               int64  `primary:"true" must_generate:"true" db:"id"`
	TelegramId       int64  `db:"telegram_id"`
	TelegramBotToken string `db:"telegram_bot_token"`
	UserId           int64  `db:"user_id"`
	User             User   `db:"-"`
}

func (u *TelegramUser) GetUser() User {
	return u.User
}

func (u *TelegramUser) SetUser(user User) {
	u.User = user
	u.UserId = user.Id
}

func (u *TelegramUser) GetTableName() string {
	return "telegram_users"
}
