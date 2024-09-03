package trc20

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	processor_module "nexus-wallet/internal/modules/transfer/processor"
)

type processor struct {
	grpcClient *client.GrpcClient
}

func NewProcessor(grpcClient *client.GrpcClient) processor_module.TransferProcessor {
	return &processor{
		grpcClient: grpcClient,
	}
}

type SignedMessage struct {
	RawData   *core.TransactionRaw
	Signature string
	Ret       []*core.Transaction_Result
}

func (p processor) Transfer(message interface{}) (*processor_module.TransferOutput, *app_error.AppError) {
	signedMessage, isValidType := message.(SignedMessage)
	if !isValidType {
		return nil, app_error.InternalError(errors.Errorf("message must be of type SignedMessage), %T provided", message))
	}

	hexBytes, err := common.HexStringToBytes(signedMessage.Signature)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.Errorf("invalid signature provided: %s", err))
	}

	_, err = p.grpcClient.Broadcast(&core.Transaction{
		RawData:   signedMessage.RawData,
		Signature: [][]byte{hexBytes},
		Ret:       signedMessage.Ret,
	})
	if err != nil {
		return nil, app_error.InvalidDataError(err)
	}

	return &processor_module.TransferOutput{Hash: ""}, nil
}
