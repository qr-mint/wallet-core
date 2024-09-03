package provider

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
)

type Limits struct {
	Min float64
	Max float64
}

type RedirectProvider interface {
	ProvideRedirectUrl(
		coinName app_enum.CoinName,
		fiatCurrency app_enum.Currency,
		amount float64,
		addressTo string,
	) (string, *app_error.AppError)
	GetLimits(coinName app_enum.CoinName, fiatCurrency app_enum.Currency) (*Limits, *app_error.AppError)
}
