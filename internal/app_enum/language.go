package app_enum

import (
	"nexus-wallet/internal/app_enum/utils"
)

type Language string

const (
	RuLanguage  Language = "ru"
	EngLanguage Language = "en"
)

func ToLanguage(value string) *Language {
	if !utils.AssertInArray(value, []string{string(RuLanguage), string(EngLanguage)}) {
		return nil
	}

	language := Language(value)
	return &language
}
