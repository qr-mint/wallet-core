package coin

import (
	"nexus-wallet/internal/app_enum"
)

type GetCryptoPriceInFiatInput struct {
	Network      app_enum.Network
	CoinName     app_enum.CoinName
	FiatCurrency app_enum.Currency
	Amount       float64
}

type GetCryptoPriceInFiatOutput struct {
	PayableAmount float64
	Price         float64
}

type GetFiatPriceInCryptoInput struct {
	Network      app_enum.Network
	CoinName     app_enum.CoinName
	FiatCurrency app_enum.Currency
	Amount       float64
}

type GetFiatPriceInCryptoOutput struct {
	PayableAmount float64
	Price         float64
}
