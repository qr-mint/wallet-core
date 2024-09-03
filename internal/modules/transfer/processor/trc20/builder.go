package trc20

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"math/big"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/app_util"
	processor_module "nexus-wallet/internal/modules/transfer/processor"
)

type builder struct {
	grpcClient *client.GrpcClient
}

func NewBuilder(grpcClient *client.GrpcClient) processor_module.TransferMessageBuilder {
	return &builder{grpcClient: grpcClient}
}

type BuiltMessage struct {
	RawData *core.TransactionRaw
	Ret     []*core.Transaction_Result
	Hash    string
}

func (b builder) BuildTokenMessage(input processor_module.BuildTokenMessageInput) (*processor_module.BuildMessageOutput, *app_error.AppError) {
	amount, err := app_util.AmountToInt(app_enum.Trc20Network, input.Amount)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not convert amount to int in trc transfer builder: %s", err))
	}

	var message BuiltMessage
	switch input.CoinName {
	case app_enum.TetherCoinName:
		result, err := b.grpcClient.TRC20Send(input.FromAddress, input.ToAddress, input.ContractAddress, big.NewInt(amount), 100000000)
		if err != nil {
			return nil, app_error.InternalError(errors.Errorf("can not TRC20Send trc transfer builder: %s", err))
		}

		hash, hashCreationErr := b.createTransactionDataHash(result.Transaction)
		if hashCreationErr != nil {
			return nil, hashCreationErr
		}
		message = BuiltMessage{RawData: result.Transaction.RawData, Ret: result.Transaction.Ret, Hash: hash}
	default:
		return nil, app_error.InvalidDataError(errors.Errorf("unsupportable coin name provided in trc transfer builder: %s", input.CoinName))
	}

	return &processor_module.BuildMessageOutput{Message: &message}, nil
}

func (b builder) BuildMessage(input processor_module.BuildMessageInput) (*processor_module.BuildMessageOutput, *app_error.AppError) {
	amount, err := app_util.AmountToInt(app_enum.Trc20Network, input.Amount)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not convert amount to int in trc transfer builder: %s", err))
	}
	result, err := b.grpcClient.Transfer(input.FromAddress, input.ToAddress, amount)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not invoke transfer grpc client method in trc transfer builder: %s", err))
	}

	hash, hashCreationErr := b.createTransactionDataHash(result.Transaction)
	if hashCreationErr != nil {
		return nil, hashCreationErr
	}
	message := BuiltMessage{RawData: result.Transaction.RawData, Ret: result.Transaction.Ret, Hash: hash}

	return &processor_module.BuildMessageOutput{Message: &message}, nil
}

func (builder) createTransactionDataHash(transaction *core.Transaction) (string, *app_error.AppError) {
	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {
		return "", app_error.InternalError(errors.Errorf("can not marshal trc20 transfer message: %s", err))
	}
	h256h := sha256.New()
	h256h.Write(rawData)

	return hex.EncodeToString(h256h.Sum(nil)), nil
}
