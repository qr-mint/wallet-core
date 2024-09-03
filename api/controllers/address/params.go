package address

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"nexus-wallet/api/utils"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/address"
	_import "nexus-wallet/internal/modules/address/import"
	"strconv"
)

type ImportRequest struct {
	Addresses []struct {
		Network string `json:"network"`
		Address string `json:"address"`
	} `json:"addresses"`
	Name         string `json:"name"`
	MnemonicHash string `json:"mnemonic_hash"`
}

func (ImportRequest) createInputFromRequest(context *gin.Context) (*address.ImportInput, *app_error.AppError) {
	var body ImportRequest
	if err := context.BindJSON(&body); err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid json data"))
	}
	var addresses []_import.AddressData
	for _, addressData := range body.Addresses {
		if network := app_enum.ToNetwork(addressData.Network); network == nil {
			return nil, app_error.InvalidDataError(fmt.Errorf("invalid network provided: %s", addressData.Network))
		}
		addresses = append(addresses, _import.AddressData{Network: app_enum.Network(addressData.Network), Address: addressData.Address})
	}

	if len(body.Name) > 20 {
		return nil, app_error.InvalidDataError(errors.New("`name` length too long"))
	}

	return &address.ImportInput{
		ImportData: _import.ImportData{
			Addresses:    addresses,
			Name:         body.Name,
			MnemonicHash: body.MnemonicHash,
		},
		UserId: context.GetInt64("userId"),
	}, nil
}

type InfoRequest struct {
}

func (InfoRequest) createInputFromRequest(context *gin.Context) (*address.GetAggregatedInfoInput, *app_error.AppError) {
	mnemonicId := context.GetInt64("mnemonicId")
	currency := app_enum.ToCurrency(context.GetHeader("Currency-Code"))
	if currency == nil {
		return nil, app_error.InvalidDataError(errors.New("Invalid currency provided. Param is required and must be valid currency."))
	}

	return &address.GetAggregatedInfoInput{
		MnemonicId: mnemonicId,
		Currency:   *currency,
	}, nil
}

type InfoResponseItem struct {
	Id                     int64             `json:"id"`
	FiatAmount             string            `json:"fiat_amount"`
	Currency               app_enum.Currency `json:"fiat_currency"`
	Amount                 string            `json:"amount"`
	Symbol                 string            `json:"symbol"`
	ImageSource            string            `json:"image_source"`
	Network                app_enum.Network  `json:"network"`
	Name                   app_enum.CoinName `json:"name"`
	DailyPriceDeltaPercent string            `json:"daily_price_delta_percent"`
}

func (InfoResponseItem) fillFromOutput(item address.InfoOutputItem) InfoResponseItem {
	return InfoResponseItem{
		Id:                     item.Id,
		FiatAmount:             utils.FormatFloatFiat(item.FiatAmount),
		Amount:                 utils.FormatFloat(item.Amount),
		Symbol:                 item.Symbol,
		ImageSource:            item.ImageSource,
		Network:                item.Network,
		Name:                   item.Name,
		Currency:               item.Currency,
		DailyPriceDeltaPercent: utils.FormatFloatFiat(item.DailyPriceDeltaPercent),
	}
}

type InfoResponse struct {
	GrowPercent string             `json:"daily_price_delta_percent"`
	FiatAmount  string             `json:"fiat_amount"`
	Currency    app_enum.Currency  `json:"fiat_currency"`
	Items       []InfoResponseItem `json:"items"`
	Name        string             `json:"name"`
}

func (InfoResponse) fillFromOutput(output address.GetAffrefatedInfoOutput) InfoResponse {
	var coins []InfoResponseItem
	for _, item := range output.Items {
		coins = append(coins, InfoResponseItem{}.fillFromOutput(item))
	}

	return InfoResponse{
		GrowPercent: utils.FormatFloatFiat(output.DailyPriceDeltaPercent),
		FiatAmount:  utils.FormatFloatFiat(output.FiatAmount),
		Currency:    output.Currency,
		Items:       coins,
		Name:        output.Name,
	}
}

type SwitchVisibilityRequest struct {
}

func (SwitchVisibilityRequest) createInputFromRequest(context *gin.Context) (*address.SwitchCoinVisibilityInput, *app_error.AppError) {
	mnemonicId := context.GetInt64("mnemonicId")
	addressCoinId, coinInInvalid := strconv.Atoi(context.Param("id"))
	if coinInInvalid != nil {
		return nil, app_error.InvalidDataError(errors.New("id must be provided and valid numeric value"))
	}

	return &address.SwitchCoinVisibilityInput{AddressCoinId: int64(addressCoinId), MnemonicId: mnemonicId}, nil
}

type CoinRequest struct {
}

func (CoinRequest) createInputFromRequest(context *gin.Context) (*address.GetCoinInput, *app_error.AppError) {
	addressCoinId, coinInInvalid := strconv.Atoi(context.Param("id"))
	if coinInInvalid != nil {
		return nil, app_error.InvalidDataError(errors.New("id must be provided and valid numeric value"))
	}
	mnemonicId := context.GetInt64("mnemonicId")

	return &address.GetCoinInput{Id: int64(addressCoinId), MnemonicId: mnemonicId}, nil
}

type CoinResponse struct {
	Id           int64             `json:"id"`
	CoinId       int64             `json:"coin_id"`
	Network      app_enum.Network  `json:"network"`
	Name         app_enum.CoinName `json:"name"`
	Symbol       string            `json:"symbol"`
	Caption      string            `json:"caption"`
	ImageSource  string            `json:"image_source"`
	Amount       string            `json:"amount"`
	Address      string            `json:"address"`
	ExplorerLink string            `json:"explorer_link"`
}

func (CoinResponse) fillFromOutput(output address.GetCoinOutput) CoinResponse {
	return CoinResponse{
		Id:           output.Id,
		CoinId:       output.CoinId,
		Network:      output.Network,
		Name:         output.Name,
		Symbol:       output.Symbol,
		Caption:      output.Caption,
		ImageSource:  output.ImageSource,
		Amount:       utils.FormatFloat(output.Amount),
		Address:      output.Address,
		ExplorerLink: output.ExplorerLink,
	}
}

type ListCoinsRequest struct {
}

func (ListCoinsRequest) createInputFromRequest(context *gin.Context) address.GetCoinsListInput {
	return address.GetCoinsListInput{
		MnemonicId:  context.GetInt64("mnemonicId"),
		OnlyVisible: context.Request.URL.Query().Get("only_visible") == "true",
	}
}

type ListCoinsResponseItem struct {
	Id          int64             `json:"id"`
	CoinId      int64             `json:"coin_id"`
	Network     app_enum.Network  `json:"network"`
	Name        app_enum.CoinName `json:"name"`
	Symbol      string            `json:"symbol"`
	Caption     string            `json:"caption"`
	ImageSource string            `json:"image_source"`
	Amount      string            `json:"amount"`
	IsVisible   bool              `json:"is_visible"`
	Address     string            `json:"address"`
}

func (ListCoinsResponseItem) fillFromOutput(output address.GetCoinListOutputItem) ListCoinsResponseItem {
	return ListCoinsResponseItem{
		Id:          output.Id,
		CoinId:      output.CoinId,
		Network:     output.Network,
		Name:        output.Name,
		Symbol:      output.Symbol,
		Caption:     output.Caption,
		ImageSource: output.ImageSource,
		IsVisible:   output.IsVisible,
		Amount:      utils.FormatFloat(output.Amount),
		Address:     output.Address,
	}
}

type ListCoinsResponse struct {
	Items []ListCoinsResponseItem `json:"items"`
}

func (ListCoinsResponse) fillFromOutput(output address.GetCoinsListOutput) ListCoinsResponse {
	var coins []ListCoinsResponseItem
	for _, item := range output.Items {
		coins = append(coins, ListCoinsResponseItem{}.fillFromOutput(item))
	}

	return ListCoinsResponse{Items: coins}
}
