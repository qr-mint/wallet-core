package trc20

import (
	"fmt"
	"gitlab.com/golib4/coingecko-client/coingecko"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/coin/provider"
	"time"
)

type tetherProvider struct {
	coingeckoClient *coingecko.Client
}

func NewTetherProvider(coingeckoClient *coingecko.Client) provider.Provider {
	return &tetherProvider{coingeckoClient: coingeckoClient}
}

func (p tetherProvider) GetPrice(currencyCode app_enum.Currency, date time.Time) (float64, error) {
	tetherPrice, err := p.coingeckoClient.GetTetherPrice(&coingecko.PriceRequest{Date: date.Format("02-01-2006")})
	if err != nil {
		return 0, fmt.Errorf("can not get tether price data: %s", err)
	}

	switch currencyCode {
	case app_enum.UsdCurrency:
		return tetherPrice.MarketData.CurrentPrice.USD, nil
	case app_enum.RubCurrency:
		return tetherPrice.MarketData.CurrentPrice.RUB, nil
	case app_enum.EurCurrency:
		return tetherPrice.MarketData.CurrentPrice.EUR, nil
	default:
		return 0, fmt.Errorf("unknown currency provided: %s", currencyCode)
	}
}

func (p tetherProvider) IsSupports(network app_enum.Network, name app_enum.CoinName) bool {
	return network == app_enum.Trc20Network && name == app_enum.TetherCoinName
}
