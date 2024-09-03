package deposit

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/deposit/enum"
)

type GetLimitsInput struct {
	MnemonicId    int64
	AddressCoinId int64
	FiatCurrency  app_enum.Currency
	ProviderName  enum.ProviderName
}

type GetLimitsOutput struct {
	Min float64
	Max float64
}

type ProvideRedirectUrlInput struct {
	MnemonicId    int64
	Amount        float64
	AddressCoinId int64
	FiatCurrency  app_enum.Currency
	ProviderName  enum.ProviderName
}

type ProvideRedirectUrlOutput struct {
	Url string
}
