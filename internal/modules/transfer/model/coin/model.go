package coin

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Coin struct {
	Id      int64             `primary:"true" must_generate:"true" db:"id"`
	Network app_enum.Network  `db:"network"`
	Name    app_enum.CoinName `db:"name"`
	IsToken bool              `db:"is_token"`
	Address *string           `db:"address"`
}

func (Coin) GetTableName() string {
	return "coins"
}

func (Coin) Clear() repository.Model {
	return &Coin{}
}

type Price struct {
	Id           int64  `primary:"true" must_generate:"true" db:"id"`
	Price        int64  `db:"price"`
	CoinId       int64  `db:"coin_id"`
	FiatCurrency string `db:"fiat_currency"`
}

func (p Price) GetTableName() string {
	return "coin_prices"
}

func (Price) Clear() repository.Model {
	return &Price{}
}
