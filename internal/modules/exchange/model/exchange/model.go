package exchange

import "nexus-wallet/pkg/repository"

type Exchange struct {
	Id          int64  `primary:"true" must_generate:"true" db:"id"`
	ExternalId  string `db:"external_id"`
	SupportLink string `db:"support_link"`
	CoinFromId  int64  `db:"address_coin_id_from"`
	CoinToId    int64  `db:"address_coin_id_to"`
	MnemonicId  int64  `db:"mnemonic_id"`
}

func (Exchange) GetTableName() string {
	return "exchanges"
}

func (Exchange) Clear() repository.Model {
	return &Exchange{}
}
