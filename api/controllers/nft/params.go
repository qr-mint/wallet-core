package nft

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"io"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/nft"
	"nexus-wallet/internal/shared/ton_message"
	"strconv"
)

type ListRequest struct {
}

func (ListRequest) createInputFromRequest(context *gin.Context) (*nft.ListInput, *app_error.AppError) {
	mnemonicId := context.GetInt64("mnemonicId")
	offset, offsetIsInvalid := strconv.Atoi(context.Request.URL.Query().Get("offset"))
	if offsetIsInvalid != nil {
		return nil, app_error.InvalidDataError(errors.New("`offset` must be provided and valid numeric value"))
	}
	limit, limitIsInvalid := strconv.Atoi(context.Request.URL.Query().Get("limit"))
	if limitIsInvalid != nil {
		return nil, app_error.InvalidDataError(errors.New("`limit` must be provided and valid numeric value"))
	}

	return &nft.ListInput{
		MnemonicId: mnemonicId,
		Limit:      uint(limit),
		Offset:     uint(offset),
	}, nil
}

type ListResponseItem struct {
	Id                    int64            `json:"id"`
	Address               string           `json:"address"`
	Name                  string           `json:"name"`
	Price                 int64            `json:"price"`
	TokenSymbol           string           `json:"token_symbol"`
	Index                 int64            `json:"index"`
	Network               app_enum.Network `json:"network"`
	CollectionAddress     string           `json:"collection_address"`
	CollectionName        string           `json:"collection_name"`
	CollectionDescription string           `json:"collection_description"`
	PreviewUrls           interface{}      `json:"preview_urls"`
}

func (ListResponseItem) fillFromModel(outputItem nft.ListOutputItem) ListResponseItem {
	return ListResponseItem{
		Id:                    outputItem.Id,
		Address:               outputItem.Address,
		Name:                  outputItem.Name,
		Price:                 outputItem.Price,
		TokenSymbol:           outputItem.TokenSymbol,
		Index:                 outputItem.Index,
		Network:               outputItem.Network,
		CollectionAddress:     outputItem.CollectionAddress,
		CollectionName:        outputItem.CollectionName,
		CollectionDescription: outputItem.CollectionDescription,
		PreviewUrls:           outputItem.PreviewUrls,
	}
}

type ListResponse struct {
	Items []ListResponseItem `json:"items"`
}

func (ListResponse) fillFromOutput(output nft.ListOutput) ListResponse {
	var responseItems []ListResponseItem
	for _, nftData := range output.Items {
		responseItems = append(responseItems, ListResponseItem{}.fillFromModel(nftData))
	}

	return ListResponse{Items: responseItems}
}

type GetRequest struct {
}

func (GetRequest) createInputFromRequest(context *gin.Context) (*nft.GetInput, *app_error.AppError) {
	mnemonicId := context.GetInt64("mnemonicId")
	id, idIsInvalid := strconv.Atoi(context.Param("id"))
	if idIsInvalid != nil {
		return nil, app_error.InvalidDataError(errors.New("`id` must be numeric value"))
	}

	return &nft.GetInput{
		MnemonicId: mnemonicId,
		NftId:      int64(id),
	}, nil
}

type GetResponse struct {
	Id                    int64            `json:"id"`
	Address               string           `json:"address"`
	Name                  string           `json:"name"`
	Price                 int64            `json:"price"`
	TokenSymbol           string           `json:"token_symbol"`
	Index                 int64            `json:"index"`
	Network               app_enum.Network `json:"network"`
	CollectionAddress     string           `json:"collection_address"`
	CollectionName        string           `json:"collection_name"`
	CollectionDescription string           `json:"collection_description"`
	PreviewUrls           interface{}      `json:"preview_urls"`
}

func (GetResponse) fillFromModel(outputItem nft.GetOutput) GetResponse {
	return GetResponse{
		Id:                    outputItem.Id,
		Address:               outputItem.Address,
		Name:                  outputItem.Name,
		Price:                 outputItem.Price,
		TokenSymbol:           outputItem.TokenSymbol,
		Index:                 outputItem.Index,
		Network:               outputItem.Network,
		CollectionAddress:     outputItem.CollectionAddress,
		CollectionName:        outputItem.CollectionName,
		CollectionDescription: outputItem.CollectionDescription,
		PreviewUrls:           outputItem.PreviewUrls,
	}
}

type GetMessageRequest struct {
}

func (GetMessageRequest) createInputFromRequest(context *gin.Context) (*nft.BuildSendMessageInput, *app_error.AppError) {
	mnemonicId := context.GetInt64("mnemonicId")
	network := app_enum.ToNetwork(context.Request.URL.Query().Get("network"))
	if network == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `network` provided. param is required and must be valid network"))
	}
	addressTo := context.Request.URL.Query().Get("address_to")
	if addressTo == "" {
		return nil, app_error.InvalidDataError(errors.New("param `address_to` is required"))
	}
	version := context.Request.URL.Query().Get("version")
	intVersion := 0
	var err error
	if version != "" {
		intVersion, err = strconv.Atoi(version)
		if err != nil {
			return nil, app_error.InvalidDataError(errors.New("invalid `version` given. param is required and must be numeric"))
		}
	}
	nftId, err := strconv.Atoi(context.Request.URL.Query().Get("nft_id"))
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `nft_id` given. param is required and must be numeric"))
	}

	return &nft.BuildSendMessageInput{
		Network:    *network,
		MnemonicId: mnemonicId,
		ToAddress:  addressTo,
		NftId:      int64(nftId),
		PublicKey:  context.Request.URL.Query().Get("public_key"),
		Version:    uint8(intVersion),
	}, nil
}

type GetMessageResponse struct {
}

type GetTonMessageResponse struct {
	DestinationAddress string         `json:"destination_address"`
	StateInit          *tlb.StateInit `json:"state_init"`
	Body               *cell.Cell     `json:"body"`
}

func (GetMessageResponse) fillFromMessage(message interface{}) (interface{}, error) {
	switch message := message.(type) {
	case *ton_message.BuiltMessage:
		return GetTonMessageResponse{DestinationAddress: message.DestinationAddress, StateInit: message.StateInit, Body: message.Body}, nil
	default:
		return nil, fmt.Errorf("unsupported message type %T", message)
	}
}

type SendRequest struct {
	Network string `json:"network"`
	NftId   int64  `json:"nft_id"`
}

func (SendRequest) createInputFromRequest(context *gin.Context) (*nft.SendInput, *app_error.AppError) {
	requestData, err := io.ReadAll(context.Request.Body)
	defer context.Request.Body.Close()
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("can not read request. invalid request"))
	}

	body := SendRequest{}
	if err := json.Unmarshal(requestData, &body); err != nil {
		return nil, app_error.InvalidDataError(errors.New("can not parse request. invalid request"))
	}

	network := app_enum.ToNetwork(body.Network)
	if network == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `network` provided. param is required and must be valid network"))
	}

	var message interface{}
	switch *network {
	case app_enum.TonNetwork:
		requestMessage := SendNftTonMessageRequest{}
		if err := json.Unmarshal(requestData, &requestMessage); err != nil {
			return nil, app_error.InvalidDataError(errors.New("can not parse ton request. invalid request"))
		}
		message = ton_message.SignedMessage{
			DestinationAddress: requestMessage.Message.DestinationAddress,
			StateInit:          requestMessage.Message.StateInit,
			Body:               requestMessage.Message.Body,
			Signature:          requestMessage.Message.Signature,
		}
	default:
		return nil, app_error.InternalError(errors.Errorf("unsupported network given:%s", *network))
	}

	return &nft.SendInput{
		Network:    *network,
		NftId:      body.NftId,
		MnemonicId: context.GetInt64("mnemonicId"),
		Message:    message,
	}, nil
}

type SendNftTonMessageRequest struct {
	Message struct {
		DestinationAddress string         `json:"destination_address"`
		StateInit          *tlb.StateInit `json:"state_init"`
		Body               *cell.Cell     `json:"body"`
		Signature          string         `json:"signature"`
	} `json:"message"`
}

type SendResponse struct {
	Hash string `json:"hash"`
}
