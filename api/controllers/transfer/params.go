package transfer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/gin-gonic/gin"
	errors2 "github.com/pkg/errors"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"io"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/transfer"
	"nexus-wallet/internal/modules/transfer/processor/trc20"
	"nexus-wallet/internal/shared/ton_message"
	"strconv"
)

type GetMessageRequest struct {
}

func (GetMessageRequest) createInputFromRequest(context *gin.Context) (*transfer.BuildTransferInput, *app_error.AppError) {
	mnemonicId := context.GetInt64("mnemonicId")
	coinName := app_enum.ToCoinName(context.Request.URL.Query().Get("coin_name"))
	if coinName == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `coin_name` provided. param is required and must be valid coin name"))
	}
	network := app_enum.ToNetwork(context.Request.URL.Query().Get("network"))
	if network == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `network` provided. param is required and must be valid network"))
	}
	addressTo := context.Request.URL.Query().Get("address_to")
	if addressTo == "" {
		return nil, app_error.InvalidDataError(errors.New("param `address_to` is required"))
	}
	amount, err := strconv.ParseFloat(context.Request.URL.Query().Get("amount"), 64)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `amount` given. param is required and must be numeric"))
	}
	version := context.Request.URL.Query().Get("version")
	intVersion := 0
	if version != "" {
		intVersion, err = strconv.Atoi(version)
		if err != nil {
			return nil, app_error.InvalidDataError(errors.New("invalid `version` given. param is required and must be numeric"))
		}
	}

	return &transfer.BuildTransferInput{
		Network:    *network,
		CoinName:   *coinName,
		MnemonicId: mnemonicId,
		ToAddress:  addressTo,
		Amount:     amount,
		PublicKey:  context.Request.URL.Query().Get("public_key"),
		Version:    uint8(intVersion),
		Comment:    context.Request.URL.Query().Get("comment"),
	}, nil
}

type GetMessageResponse struct {
}

type GetTonMessageResponse struct {
	DestinationAddress string         `json:"destination_address"`
	StateInit          *tlb.StateInit `json:"state_init"`
	Body               *cell.Cell     `json:"body"`
}

type GetTrc20MessageResponse struct {
	RawData *core.TransactionRaw       `json:"raw_data"`
	Ret     []*core.Transaction_Result `json:"ret"`
	Hash    string                     `json:"hash"`
}

func (GetMessageResponse) fillFromMessage(message interface{}) (interface{}, error) {
	switch message := message.(type) {
	case *ton_message.BuiltMessage:
		return GetTonMessageResponse{DestinationAddress: message.DestinationAddress, StateInit: message.StateInit, Body: message.Body}, nil
	case *trc20.BuiltMessage:
		return GetTrc20MessageResponse{RawData: message.RawData, Ret: message.Ret, Hash: message.Hash}, nil
	default:
		return nil, fmt.Errorf("unsupported message type %T", message)
	}
}

type TransferRequest struct {
	Network  string  `json:"network"`
	CoinName string  `json:"coin_name"`
	Amount   float64 `json:"amount"`
}

func (TransferRequest) createInputFromRequest(context *gin.Context) (*transfer.TransferInput, *app_error.AppError) {
	mnemonicId := context.GetInt64("mnemonicId")
	requestData, err := io.ReadAll(context.Request.Body)
	defer context.Request.Body.Close()
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("can not read request. invalid request"))
	}

	body := TransferRequest{}
	if err := json.Unmarshal(requestData, &body); err != nil {
		return nil, app_error.InvalidDataError(errors.New("can not parse request. invalid request"))
	}

	network := app_enum.ToNetwork(body.Network)
	if network == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `network` provided. param is required and must be valid network"))
	}
	coinName := app_enum.ToCoinName(body.CoinName)
	if coinName == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `coin_name` provided. param is required and must be valid coin name"))
	}

	var message interface{}
	switch *network {
	case app_enum.TonNetwork:
		requestMessage := TransferTonMessageRequest{}
		if err := json.Unmarshal(requestData, &requestMessage); err != nil {
			return nil, app_error.InvalidDataError(errors.New("can not parse ton request. invalid request"))
		}
		message = ton_message.SignedMessage{
			DestinationAddress: requestMessage.Message.DestinationAddress,
			StateInit:          requestMessage.Message.StateInit,
			Body:               requestMessage.Message.Body,
			Signature:          requestMessage.Message.Signature,
		}
	case app_enum.Trc20Network:
		requestMessage := TransferTrc20MessageRequest{}
		if err := json.Unmarshal(requestData, &requestMessage); err != nil {
			return nil, app_error.InvalidDataError(errors.New("can not parse trc20 request. invalid request"))
		}
		message = trc20.SignedMessage{
			RawData:   requestMessage.Message.RawData,
			Ret:       requestMessage.Message.Ret,
			Signature: requestMessage.Message.Signature,
		}
	default:
		return nil, app_error.InternalError(errors2.Errorf("unsupported network given:%s", *network))
	}

	return &transfer.TransferInput{
		Network:    *network,
		CoinName:   *coinName,
		Amount:     body.Amount,
		MnemonicId: mnemonicId,
		Message:    message,
	}, nil
}

type TransferTonMessageRequest struct {
	Message struct {
		DestinationAddress string         `json:"destination_address"`
		StateInit          *tlb.StateInit `json:"state_init"`
		Body               *cell.Cell     `json:"body"`
		Signature          string         `json:"signature"`
	} `json:"message"`
}

type TransferTrc20MessageRequest struct {
	Message struct {
		RawData   *core.TransactionRaw       `json:"raw_data"`
		Ret       []*core.Transaction_Result `json:"ret"`
		Signature string                     `json:"signature"`
	} `json:"message"`
}

type TransferResponse struct {
	Hash string `json:"hash"`
}
