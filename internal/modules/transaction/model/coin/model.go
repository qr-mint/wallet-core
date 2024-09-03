package coin

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Coin struct {
	Id      int64             `primary:"true" must_generate:"true" db:"id"`
	Network app_enum.Network  `db:"network"`
	Name    app_enum.CoinName `db:"name"`
	Address *string           `db:"address"`
}

func (Coin) GetTableName() string {
	return "coins"
}

func (Coin) Clear() repository.Model {
	return &Coin{}
}
