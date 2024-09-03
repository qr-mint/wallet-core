package coin

import "nexus-wallet/internal/app_enum"

type Coin struct {
	Id   int64             `primary:"true" must_generate:"true" db:"id"`
	Name app_enum.CoinName `db:"name"`
}

func (Coin) GetTableName() string {
	return "coins"
}
