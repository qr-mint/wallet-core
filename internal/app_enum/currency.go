package app_enum

import (
	"nexus-wallet/internal/app_enum/utils"
)

type Currency string

const (
	UsdCurrency Currency = "usd"
	RubCurrency Currency = "rub"
	EurCurrency Currency = "eur"
)

func ToCurrency(value string) *Currency {
	if !utils.AssertInArray(value, []string{string(UsdCurrency), string(RubCurrency), string(EurCurrency)}) {
		return nil
	}

	currency := Currency(value)
	return &currency
}

func GetCurrencies() []Currency {
	return []Currency{
		UsdCurrency,
		RubCurrency,
		EurCurrency,
	}
}
