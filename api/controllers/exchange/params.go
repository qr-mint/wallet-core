package exchange

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"nexus-wallet/api/utils"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/exchange"
	"strconv"
)

type ListRequest struct {
	MnemonicId int64
	Limit      uint
	Offset     uint
}

func (ListRequest) createInputFromRequest(context *gin.Context) (*exchange.ListInput, *app_error.AppError) {
	offset, offsetIsInvalid := strconv.Atoi(context.Request.URL.Query().Get("offset"))
	if offsetIsInvalid != nil {
		return nil, app_error.InvalidDataError(errors.New("`offset` must be provided and valid numeric value"))
	}
	limit, limitIsInvalid := strconv.Atoi(context.Request.URL.Query().Get("limit"))
	if limitIsInvalid != nil {
		return nil, app_error.InvalidDataError(errors.New("`limit` must be provided and valid numeric value"))
	}

	mnemonicId := context.GetInt64("mnemonicId")

	return &exchange.ListInput{
		MnemonicId: mnemonicId,
		Limit:      uint(limit),
		Offset:     uint(offset),
	}, nil
}

type ListResponseItem struct {
	ExternalId   string            `json:"external_id"`
	SupportLink  string            `json:"support_link"`
	CoinFromName app_enum.CoinName `json:"coin_from_name"`
	CoinToName   app_enum.CoinName `json:"coin_to_name"`
}

func (ListResponseItem) fillFromOutput(output exchange.ListOutputItem) ListResponseItem {
	return ListResponseItem{
		ExternalId:   output.ExternalId,
		SupportLink:  output.SupportLink,
		CoinFromName: output.CoinFromName,
		CoinToName:   output.CoinToName,
	}
}

type ListResponse struct {
	Items []ListResponseItem `json:"items"`
}

func (ListResponse) fillFromOutput(output exchange.ListOutput) ListResponse {
	var items []ListResponseItem
	for _, outputItem := range output.Items {
		items = append(items, ListResponseItem{}.fillFromOutput(outputItem))
	}

	return ListResponse{Items: items}
}

type GetAddressForTransferRequest struct {
}

func (GetAddressForTransferRequest) createInputFromRequest(context *gin.Context) (*exchange.ProvideAddressForTransferInput, *app_error.AppError) {
	amount, err := strconv.ParseFloat(context.Request.URL.Query().Get("amount"), 64)
	if err != nil || amount <= 0 {
		return nil, app_error.InvalidDataError(errors.New("invalid amount given"))
	}
	addressCoinIdFrom, err := strconv.Atoi(context.Request.URL.Query().Get("address_coin_id_from"))
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `address_coin_id_from` given"))
	}
	addressCoinIdTo, err := strconv.Atoi(context.Request.URL.Query().Get("address_coin_id_to"))
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `address_coin_id_to` given"))
	}
	if addressCoinIdFrom == addressCoinIdTo {
		return nil, app_error.InvalidDataError(errors.New("`address_coin_id_from` and `address_coin_id_to` can not be equal"))
	}

	mnemonicId := context.GetInt64("mnemonicId")

	return &exchange.ProvideAddressForTransferInput{
		AddressCoinIdFrom: int64(addressCoinIdFrom),
		AddressCoinIdTo:   int64(addressCoinIdTo),
		MnemonicId:        mnemonicId,
		Amount:            amount,
	}, nil
}

type GetAddressForTransferResponse struct {
	PayInAddress  string `json:"pay_in_address"`
	TransactionId string `json:"transaction_id"`
}

func (GetAddressForTransferResponse) fillFromOutput(output exchange.ProvideAddressForTransferOutput) GetAddressForTransferResponse {
	return GetAddressForTransferResponse{
		PayInAddress:  output.PayInAddress,
		TransactionId: output.TransactionId,
	}
}

type GetExchangeAmountRequest struct {
}

func (GetExchangeAmountRequest) createInputFromRequest(context *gin.Context) (*exchange.GetExchangeAmountInput, *app_error.AppError) {
	amount, err := strconv.ParseFloat(context.Request.URL.Query().Get("send_amount"), 64)
	if err != nil || amount <= 0 {
		return nil, app_error.InvalidDataError(errors.New("invalid `send_amount` given"))
	}
	addressCoinIdFrom, err := strconv.Atoi(context.Request.URL.Query().Get("address_coin_id_from"))
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `address_coin_id_from` given"))
	}
	addressCoinIdTo, err := strconv.Atoi(context.Request.URL.Query().Get("address_coin_id_to"))
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `address_coin_id_to` given"))
	}
	if addressCoinIdFrom == addressCoinIdTo {
		return nil, app_error.InvalidDataError(errors.New("`address_coin_id_from` and `address_coin_id_to` can not be equal"))
	}

	mnemonicId := context.GetInt64("mnemonicId")

	return &exchange.GetExchangeAmountInput{
		AddressCoinIdFrom: int64(addressCoinIdFrom),
		AddressCoinIdTo:   int64(addressCoinIdTo),
		MnemonicId:        mnemonicId,
		SendAmount:        amount,
	}, nil

}

type GetExchangeAmountResponse struct {
	ReceiveAmount string `json:"receive_amount"`
}

func (GetExchangeAmountResponse) fillFromOutput(output exchange.GetExchangeAmountOutput) GetExchangeAmountResponse {
	return GetExchangeAmountResponse{
		ReceiveAmount: utils.FormatFloat(output.ReceiveAmount),
	}
}

type GetLimitsRequest struct {
}

func (GetLimitsRequest) createInputFromRequest(context *gin.Context) (*exchange.GetLimitsInput, *app_error.AppError) {
	addressCoinIdFrom, err := strconv.Atoi(context.Request.URL.Query().Get("address_coin_id_from"))
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `address_coin_id_from` given"))
	}
	addressCoinIdTo, err := strconv.Atoi(context.Request.URL.Query().Get("address_coin_id_to"))
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `address_coin_id_to` given"))
	}
	mnemonicId := context.GetInt64("mnemonicId")

	if addressCoinIdFrom == addressCoinIdTo {
		return nil, app_error.InvalidDataError(errors.New("`address_coin_id_from` and `address_coin_id_to` can not be equal"))
	}

	return &exchange.GetLimitsInput{
		MnemonicId:        mnemonicId,
		AddressCoinIdFrom: int64(addressCoinIdFrom),
		AddressCoinIdTo:   int64(addressCoinIdTo),
	}, nil
}

type GetLimitsResponse struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

func (GetLimitsResponse) fillFromOutput(output exchange.GetLimitsOutput) GetLimitsResponse {
	return GetLimitsResponse{
		Min: fmt.Sprintf("%.2f", output.Min),
		Max: fmt.Sprintf("%.2f", output.Max),
	}
}
