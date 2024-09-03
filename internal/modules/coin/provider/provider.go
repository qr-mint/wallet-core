package provider

import (
	"nexus-wallet/internal/app_enum"
	"time"
)

type Provider interface {
	GetPrice(currencyCode app_enum.Currency, date time.Time) (float64, error)
	IsSupports(network app_enum.Network, name app_enum.CoinName) bool
}
