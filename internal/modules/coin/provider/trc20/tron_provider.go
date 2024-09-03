package trc20

import (
	"fmt"
	"gitlab.com/golib4/coingecko-client/coingecko"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/coin/provider"
	"time"
)

type tronProvider struct {
	coingeckoClient *coingecko.Client
}

func NewTronProvider(coingeckoClient *coingecko.Client) provider.Provider {
	return &tronProvider{coingeckoClient: coingeckoClient}
}

func (p tronProvider) GetPrice(currencyCode app_enum.Currency, date time.Time) (float64, error) {
	tronPrice, err := p.coingeckoClient.GetTronPrice(&coingecko.PriceRequest{Date: date.Format("02-01-2006")})
	if err != nil {
		return 0, fmt.Errorf("can not get tron price data: %s", err)
	}

	switch currencyCode {
	case app_enum.UsdCurrency:
		return tronPrice.MarketData.CurrentPrice.USD, nil
	case app_enum.RubCurrency:
		return tronPrice.MarketData.CurrentPrice.RUB, nil
	case app_enum.EurCurrency:
		return tronPrice.MarketData.CurrentPrice.EUR, nil
	default:
		return 0, fmt.Errorf("unknown currency provided: %s", currencyCode)
	}
}

func (p tronProvider) IsSupports(network app_enum.Network, name app_enum.CoinName) bool {
	return network == app_enum.Trc20Network && name == app_enum.TronCoinName
}
