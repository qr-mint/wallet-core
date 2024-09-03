package address

import (
	"fmt"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/internal/modules/address/import"
	"nexus-wallet/internal/modules/address/model/address"
	"nexus-wallet/internal/modules/address/model/coin"
	"nexus-wallet/internal/modules/address/util"
)

type ImportInput struct {
	ImportData _import.ImportData
	UserId     int64
}

type SwitchCoinVisibilityInput struct {
	AddressCoinId int64
	MnemonicId    int64
}

type GetCoinInput struct {
	Id         int64
	MnemonicId int64
}

type GetCoinOutput struct {
	Id           int64
	Network      app_enum.Network
	Name         app_enum.CoinName
	Symbol       string
	Caption      string
	ImageSource  string
	CoinId       int64
	Address      string
	Amount       float64
	ExplorerLink string
}

func (GetCoinOutput) fillFromModel(coin coin.Coin, addressCoin address.AddressCoin, address address.Address) (*GetCoinOutput, error) {
	floatAmount, err := app_util.AmountToFloat(coin.Network, addressCoin.Amount)
	if err != nil {
		return nil, fmt.Errorf("can't convert amount to float: %w", err)
	}

	link, err := util.ProvideExplorerLink(address.Network, address.Address)
	if err != nil {
		return nil, fmt.Errorf("can't provide explorer link: %w", err)
	}

	return &GetCoinOutput{
		Id:           addressCoin.Id,
		Network:      coin.Network,
		Name:         coin.Name,
		Symbol:       coin.Symbol,
		Caption:      coin.Caption,
		ImageSource:  coin.ImageSource,
		CoinId:       coin.Id,
		Address:      address.Address,
		Amount:       floatAmount,
		ExplorerLink: link,
	}, nil
}

type GetCoinsListInput struct {
	MnemonicId  int64
	OnlyVisible bool
}

type GetCoinListOutputItem struct {
	Id          int64
	Network     app_enum.Network
	Name        app_enum.CoinName
	Symbol      string
	Caption     string
	ImageSource string
	CoinId      int64
	Amount      float64
	Address     string
	IsVisible   bool
}

func (GetCoinListOutputItem) fillFromModel(coin coin.Coin, addressCoin address.AddressCoin, address address.Address) (*GetCoinListOutputItem, error) {
	floatAmount, err := app_util.AmountToFloat(coin.Network, addressCoin.Amount)
	if err != nil {
		return nil, fmt.Errorf("can't convert amount to float: %w", err)
	}

	return &GetCoinListOutputItem{
		Id:          addressCoin.Id,
		Network:     coin.Network,
		Name:        coin.Name,
		Symbol:      coin.Symbol,
		Caption:     coin.Caption,
		ImageSource: coin.ImageSource,
		CoinId:      coin.Id,
		Amount:      floatAmount,
		Address:     address.Address,
		IsVisible:   addressCoin.IsVisible,
	}, nil
}

type GetCoinsListOutput struct {
	Items []GetCoinListOutputItem
}

type GetAggregatedInfoInput struct {
	MnemonicId int64
	Currency   app_enum.Currency
}

type InfoOutputItem struct {
	Id                     int64
	FiatAmount             float64
	Currency               app_enum.Currency
	Amount                 float64
	Symbol                 string
	ImageSource            string
	DailyPriceDeltaPercent float64
	Network                app_enum.Network
	Name                   app_enum.CoinName
}

type GetAffrefatedInfoOutput struct {
	DailyPriceDeltaPercent float64
	FiatAmount             float64
	Currency               app_enum.Currency
	Items                  []InfoOutputItem
	Name                   string
}
