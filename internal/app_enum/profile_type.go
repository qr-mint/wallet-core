package app_enum

import (
	"nexus-wallet/internal/app_enum/utils"
)

type ProfileType string

const (
	TelegramProfileType ProfileType = "telegram"
)

func ToProfileType(value string) *ProfileType {
	if !utils.AssertInArray(value, []string{string(TelegramProfileType)}) {
		return nil
	}

	profileType := ProfileType(value)
	return &profileType
}
