package transaction

import (
	"fmt"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/internal/modules/transaction/enum"
	"nexus-wallet/internal/modules/transaction/model/coin"
	"nexus-wallet/internal/modules/transaction/model/transaction"
	"nexus-wallet/internal/modules/transaction/util"
	"time"
)

type ListInput struct {
	AddressCoinId *int64
	OnlyOut       *bool
	MnemonicId    int64
	Limit         uint
	Offset        uint
}

type ListOutputItem struct {
	Id           int64
	Hash         string
	Amount       float64
	Network      app_enum.Network
	CoinName     app_enum.CoinName
	Status       enum.Status
	Type         enum.Type
	AddressTo    string
	AddressFrom  string
	CreatedAt    time.Time
	ExplorerLink string
}

func (ListOutputItem) fillFromModel(transaction transaction.BlockchainTransaction, coin coin.Coin) (*ListOutputItem, error) {
	amount, err := app_util.AmountToFloat(coin.Network, transaction.Amount)
	if err != nil {
		return nil, fmt.Errorf("can not format amount to float: %s", err)
	}
	explorerLink, err := util.ProvideExplorerLink(coin.Network, transaction.Hash)
	if err != nil {
		return nil, fmt.Errorf("can not provide explorer link for ListOutputItem: %s", err)
	}
	return &ListOutputItem{
		Id:           transaction.Id,
		Hash:         transaction.Hash,
		Amount:       amount,
		Status:       transaction.Status,
		Type:         transaction.Type,
		CoinName:     coin.Name,
		Network:      coin.Network,
		AddressFrom:  transaction.AddressFrom,
		AddressTo:    transaction.AddressTo,
		CreatedAt:    transaction.CreatedAt,
		ExplorerLink: explorerLink,
	}, nil
}

type ListOutput struct {
	Items []ListOutputItem
}

func (ListOutput) fillFromModel(transactions []*transaction.BlockchainTransaction, coinsMappedByIds map[int64]coin.Coin) (*ListOutput, error) {
	var outputItems []ListOutputItem
	for _, transactionData := range transactions {
		outputItem, err := ListOutputItem{}.fillFromModel(*transactionData, coinsMappedByIds[transactionData.CoinId])
		if err != nil {
			return nil, fmt.Errorf("can not fill ListOutputItem from model: %s", err)
		}

		outputItems = append(outputItems, *outputItem)
	}

	return &ListOutput{Items: outputItems}, nil
}
