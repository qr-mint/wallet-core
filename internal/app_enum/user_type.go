package app_enum

import (
	"nexus-wallet/internal/app_enum/utils"
)

type UserType string

const (
	TelegramUserType UserType = "telegram"
)

func ToUserType(value string) *UserType {
	if !utils.AssertInArray(value, []string{string(TelegramUserType)}) {
		return nil
	}

	userType := UserType(value)
	return &userType
}
