package address

type Address struct {
	Id      int64  `primary:"true" must_generate:"true" db:"id"`
	Address string `db:"address"`
}

func (Address) GetTableName() string {
	return "wallet_addresses"
}

type AddressCoin struct {
	Id     int64 `primary:"true" must_generate:"true" db:"id"`
	Amount int64 `db:"amount"`
}

func (AddressCoin) GetTableName() string {
	return "wallet_address_coins"
}
