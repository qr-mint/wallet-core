package exchange

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/exchange/model/coin"
	"nexus-wallet/internal/modules/exchange/model/exchange"
)

type ListInput struct {
	MnemonicId int64
	Limit      uint
	Offset     uint
}

type ListOutputItem struct {
	ExternalId   string
	SupportLink  string
	CoinFromName app_enum.CoinName
	CoinToName   app_enum.CoinName
}

func (ListOutputItem) fillFromModel(exchange exchange.Exchange, coinFrom coin.Coin, coinTo coin.Coin) ListOutputItem {
	return ListOutputItem{
		ExternalId:   exchange.ExternalId,
		SupportLink:  exchange.SupportLink,
		CoinFromName: coinFrom.Name,
		CoinToName:   coinTo.Name,
	}

}

type ListOutput struct {
	Items []ListOutputItem
}

func (ListOutput) fillFromModel(exchanges []*exchange.Exchange, coinsMappedByIds map[int64]coin.Coin) *ListOutput {
	var outputItems []ListOutputItem
	for _, exchangeDate := range exchanges {
		outputItem := ListOutputItem{}.fillFromModel(*exchangeDate, coinsMappedByIds[exchangeDate.CoinFromId], coinsMappedByIds[exchangeDate.CoinToId])
		outputItems = append(outputItems, outputItem)
	}

	return &ListOutput{Items: outputItems}
}

type ProvideAddressForTransferInput struct {
	AddressCoinIdFrom int64
	AddressCoinIdTo   int64
	MnemonicId        int64
	Amount            float64
}

type ProvideAddressForTransferOutput struct {
	PayInAddress  string
	TransactionId string
}

type GetExchangeAmountInput struct {
	AddressCoinIdFrom int64
	AddressCoinIdTo   int64
	SendAmount        float64
	MnemonicId        int64
}

type GetExchangeAmountOutput struct {
	ReceiveAmount float64
}

type GetLimitsInput struct {
	AddressCoinIdFrom int64
	AddressCoinIdTo   int64
	MnemonicId        int64
}

type GetLimitsOutput struct {
	Min float64
	Max float64
}
