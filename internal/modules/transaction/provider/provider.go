package provider

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/transaction/enum"
	"time"
)

type BlockchainTransactionData struct {
	From      string
	To        string
	Amount    int64
	Hash      string
	Type      enum.Type
	CoinName  app_enum.CoinName
	Network   app_enum.Network
	Status    enum.Status
	CreatedAt time.Time
}

type Provider interface {
	GetTransactions(account string, limit int32, timestampFrom time.Time) ([]BlockchainTransactionData, error)
}
