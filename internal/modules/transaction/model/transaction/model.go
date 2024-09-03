package transaction

import (
	"nexus-wallet/internal/modules/transaction/enum"
	"nexus-wallet/pkg/repository"
	"time"
)

type BlockchainTransaction struct {
	Id          int64       `primary:"true" must_generate:"true" db:"id"`
	Hash        string      `db:"hash"`
	Amount      int64       `db:"amount"`
	Address     string      `db:"address"`
	AddressTo   string      `db:"address_to"`
	AddressFrom string      `db:"address_from"`
	Status      enum.Status `db:"status"`
	Type        enum.Type   `db:"type"`
	CoinId      int64       `db:"coin_id"`
	CreatedAt   time.Time   `db:"created_at"`
}

func (BlockchainTransaction) GetTableName() string {
	return "blockchain_transactions"
}

func (BlockchainTransaction) Clear() repository.Model {
	return &BlockchainTransaction{}
}
