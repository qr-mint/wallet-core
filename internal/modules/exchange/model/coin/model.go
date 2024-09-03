package coin

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Price struct {
	Id     int64 `primary:"true" must_generate:"true" db:"id"`
	Price  int64 `db:"price"`
	CoinId int64 `db:"coin_id"`
}

func (p Price) GetTableName() string {
	return "coin_prices"
}

func (Price) Clear() repository.Model {
	return &Price{}
}

type Coin struct {
	Id   int64             `primary:"true" must_generate:"true" db:"id"`
	Name app_enum.CoinName `db:"name"`
}

func (Coin) GetTableName() string {
	return "coins"
}

func (Coin) Clear() repository.Model {
	return &Coin{}
}
