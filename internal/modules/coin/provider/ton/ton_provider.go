package ton

import (
	"fmt"
	"gitlab.com/golib4/tonhubapi-client/tonhubapi"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/coin/provider"
	"time"
)

type tonProvider struct {
	tonhubapiClient *tonhubapi.Client
}

func NewTonProvider(tonhubapiClient *tonhubapi.Client) provider.Provider {
	return &tonProvider{tonhubapiClient: tonhubapiClient}
}

func (p tonProvider) GetPrice(currencyCode app_enum.Currency, date time.Time) (float64, error) {
	tonPrice, err := p.tonhubapiClient.GetPrice(&tonhubapi.PriceRequest{Date: date.Format("02-01-2006")})
	if err != nil {
		return 0, fmt.Errorf("can not get ton price data: %s", err)
	}

	switch currencyCode {
	case app_enum.UsdCurrency:
		return tonPrice.Price.USD, nil
	case app_enum.RubCurrency:
		return tonPrice.Price.USD * tonPrice.Price.Rates.RUB, nil
	case app_enum.EurCurrency:
		return tonPrice.Price.USD * tonPrice.Price.Rates.EUR, nil
	default:
		return 0, fmt.Errorf("unknown currency provided: %s", currencyCode)
	}
}

func (p tonProvider) IsSupports(network app_enum.Network, name app_enum.CoinName) bool {
	return network == app_enum.TonNetwork && name == app_enum.TonCoinName
}
