package transaction

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"nexus-wallet/api/utils"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/transaction"
	"nexus-wallet/internal/modules/transaction/enum"
	"strconv"
	"time"
)

type ListRequest struct {
}

func (ListRequest) createInputFromRequest(context *gin.Context) (*transaction.ListInput, *app_error.AppError) {
	mnemonicId := context.GetInt64("mnemonicId")
	offset, offsetIsInvalid := strconv.Atoi(context.Request.URL.Query().Get("offset"))
	if offsetIsInvalid != nil {
		return nil, app_error.InvalidDataError(errors.New("`offset` must be provided and valid numeric value"))
	}
	limit, limitIsInvalid := strconv.Atoi(context.Request.URL.Query().Get("limit"))
	if limitIsInvalid != nil {
		return nil, app_error.InvalidDataError(errors.New("`limit` must be provided and valid numeric value"))
	}
	addressCoinIdParam := context.Request.URL.Query().Get("address_coin_id")
	var addressCoinId *int64
	if addressCoinIdParam != "" {
		addressCoinIdInt, err := strconv.ParseInt(addressCoinIdParam, 10, 64)
		if err != nil {
			return nil, app_error.InvalidDataError(errors.New("`address_coin_id` must be a valid numeric value"))
		}
		addressCoinId = &addressCoinIdInt
	}
	onlyOutParam := context.Request.URL.Query().Get("only_out")
	var onlyOut *bool
	if onlyOutParam == "true" {
		onlyOutValue := true
		onlyOut = &onlyOutValue
	}

	return &transaction.ListInput{
		AddressCoinId: addressCoinId,
		OnlyOut:       onlyOut,
		MnemonicId:    mnemonicId,
		Limit:         uint(limit),
		Offset:        uint(offset),
	}, nil
}

type ListResponseItem struct {
	Id           int64             `json:"id"`
	Hash         string            `json:"hash"`
	Amount       string            `json:"amount"`
	Network      app_enum.Network  `json:"network"`
	CoinName     app_enum.CoinName `json:"coin_name"`
	Status       enum.Status       `json:"status"`
	Type         enum.Type         `json:"type"`
	AddressFrom  string            `json:"address_from"`
	AddressTo    string            `json:"address_to"`
	CreatedAt    string            `json:"created_at"`
	ExplorerLink string            `json:"explorer_link"`
}

func (ListResponseItem) fillFromOutput(output transaction.ListOutputItem) ListResponseItem {
	return ListResponseItem{
		Id:           output.Id,
		Hash:         output.Hash,
		Amount:       utils.FormatFloat(output.Amount),
		Status:       output.Status,
		Type:         output.Type,
		CoinName:     output.CoinName,
		Network:      output.Network,
		AddressFrom:  output.AddressFrom,
		AddressTo:    output.AddressTo,
		CreatedAt:    output.CreatedAt.Format(time.RFC3339),
		ExplorerLink: output.ExplorerLink,
	}
}

type ListResponse struct {
	Items []ListResponseItem `json:"items"`
}

func (ListResponse) fillFromOutput(output transaction.ListOutput) ListResponse {
	var items []ListResponseItem
	for _, outputItem := range output.Items {
		items = append(items, ListResponseItem{}.fillFromOutput(outputItem))
	}

	return ListResponse{Items: items}
}
