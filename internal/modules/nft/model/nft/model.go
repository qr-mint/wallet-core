package nft

import (
	"nexus-wallet/pkg/repository"
)

type Nft struct {
	Id                    int64  `primary:"true" must_generate:"true" db:"id"`
	Address               string `db:"address"`
	Name                  string `db:"name"`
	Price                 int64  `db:"price"`
	TokenSymbol           string `db:"token_symbol"`
	Index                 int64  `db:"index"`
	CollectionAddress     string `db:"collection_address"`
	CollectionName        string `db:"collection_name"`
	CollectionDescription string `db:"collection_description"`
	PreviewData           []byte `db:"previews_urls"`
	AddressId             int64  `db:"address_id"`
}

func (Nft) GetTableName() string {
	return "wallet_address_nfts"
}

func (Nft) Clear() repository.Model {
	return &Nft{}
}
