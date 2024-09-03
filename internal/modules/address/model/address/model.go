package address

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Address struct {
	Id         int64            `primary:"true" must_generate:"true" db:"id"`
	Address    string           `db:"address"`
	Network    app_enum.Network `db:"network"`
	MnemonicId int64            `db:"mnemonic_id"`
}

func (Address) GetTableName() string {
	return "wallet_addresses"
}

func (Address) Clear() repository.Model {
	return &Address{}
}

type AddressCoin struct {
	Id        int64  `primary:"true" must_generate:"true" db:"id"`
	Amount    int64  `db:"amount"`
	IsVisible bool   `db:"is_visible"`
	Address   string `db:"address"`
	CoinId    int64  `db:"coin_id"`
	AddressId int64  `db:"address_id"`
}

func (AddressCoin) GetTableName() string {
	return "wallet_address_coins"
}

func (AddressCoin) Clear() repository.Model {
	return &AddressCoin{}
}
