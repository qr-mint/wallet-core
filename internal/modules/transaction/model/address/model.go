package address

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Address struct {
	Id      int64            `primary:"true" must_generate:"true" db:"id"`
	Address string           `db:"address"`
	Network app_enum.Network `db:"network"`
}

func (Address) GetTableName() string {
	return "wallet_addresses"
}

func (Address) Clear() repository.Model {
	return &Address{}
}
