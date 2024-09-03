package transfer

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/golib4/logger/logger"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/internal/modules/transfer/model/address"
	"nexus-wallet/internal/modules/transfer/model/coin"
	"nexus-wallet/internal/modules/transfer/processor"
	"nexus-wallet/pkg/cache"
	"time"
)

type Service struct {
	transferProcessors      map[app_enum.Network]processor.TransferProcessor
	transferMessageBuilders map[app_enum.Network]processor.TransferMessageBuilder
	coinRepository          *coin.Repository
	addressRepository       *address.Repository
	addressCoinRepository   *address.AddressCoinRepository
	cacher                  cache.Cacher
	logger                  logger.Logger
}

func NewService(
	transferProcessors map[app_enum.Network]processor.TransferProcessor,
	transferMessageBuilders map[app_enum.Network]processor.TransferMessageBuilder,
	coinRepository *coin.Repository,
	addressRepository *address.Repository,
	addressCoinRepository *address.AddressCoinRepository,
	cacher cache.Cacher,
	logger logger.Logger,
) *Service {
	return &Service{
		transferProcessors:      transferProcessors,
		transferMessageBuilders: transferMessageBuilders,
		coinRepository:          coinRepository,
		addressRepository:       addressRepository,
		addressCoinRepository:   addressCoinRepository,
		cacher:                  cacher,
		logger:                  logger,
	}
}

func (s Service) BuildTransfer(input BuildTransferInput) (*BuildTransferOutput, *app_error.AppError) {
	isValidAddress, err := app_util.IsValidAddress(input.ToAddress, input.Network)
	if !isValidAddress {
		return nil, app_error.InvalidDataError(fmt.Errorf("incorrect to address: %s", input.ToAddress))
	}
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not validate address: %s", err))
	}

	coinData, err := s.coinRepository.Find(coin.FindOptions{Network: input.Network, CoinName: input.CoinName}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not find coin: %s", err))
	}
	if coinData == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("coin of network %s, not found: %s", input.Network, input.CoinName))
	}

	fromAddress, err := s.addressRepository.Find(address.FindOptions{MnemonicId: input.MnemonicId, Network: input.Network}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not find address: %s", err))
	}
	if fromAddress == nil {
		return nil, app_error.InternalError(errors.Errorf("address of network %s, not found for mnemonic id: %d", input.Network, input.MnemonicId))
	}
	if fromAddress.Address == input.ToAddress {
		return nil, app_error.InvalidDataError(errors.New("addressTo can not be equal to addressFrom"))
	}

	transferMessageBuilder, transferMessageBuilderFounded := s.transferMessageBuilders[input.Network]
	if !transferMessageBuilderFounded {
		return nil, app_error.InternalError(errors.Errorf("not found transfer processor by network: %s", input.Network))
	}

	if !coinData.IsToken {
		toTransaction, buildTransferErr := transferMessageBuilder.BuildMessage(processor.BuildMessageInput{
			Amount:      input.Amount,
			ToAddress:   input.ToAddress,
			FromAddress: fromAddress.Address,
			Options:     processor.BuildMessageWalletInputOptions{Version: input.Version, PublicKey: input.PublicKey, Comment: input.Comment},
		})
		if buildTransferErr != nil {
			return nil, buildTransferErr
		}

		return &BuildTransferOutput{Message: toTransaction.Message}, nil
	}

	toTransaction, buildTransferErr := transferMessageBuilder.BuildTokenMessage(processor.BuildTokenMessageInput{
		Amount:          input.Amount,
		ToAddress:       input.ToAddress,
		FromAddress:     fromAddress.Address,
		CoinName:        coinData.Name,
		ContractAddress: *coinData.Address,
		Options:         processor.BuildMessageWalletInputOptions{Comment: input.Comment},
	})
	if buildTransferErr != nil {
		return nil, buildTransferErr
	}
	return &BuildTransferOutput{Message: toTransaction.Message}, nil
}

const transferKey = "last_transfer_key_mnemonic_"

type LastMessage struct {
	Message interface{} `json:"message"`
}

func (s Service) Transfer(input TransferInput) (*TransferOutput, *app_error.AppError) {
	var lastMessage LastMessage
	lastMessageKey := fmt.Sprintf("%s%d", transferKey, input.MnemonicId)
	err := s.cacher.Get(lastMessageKey, &lastMessage)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not get last message: %s", err))
	}
	if input.Message == lastMessage {
		return nil, app_error.IllegalOperationError(errors.New("same message send in short interval was avoided"))
	}

	transferProcessor, transferProcessorFounded := s.transferProcessors[input.Network]
	if !transferProcessorFounded {
		return nil, app_error.InternalError(errors.Errorf("not found transfer processor by network: %s", input.Network))
	}

	result, transferErr := transferProcessor.Transfer(input.Message)
	if transferErr != nil {
		return nil, transferErr
	}

	err = s.cacher.SetWithTTL(lastMessageKey, LastMessage{Message: input.Message}, 5*time.Second)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not set transfer last message: %s", err))
	}

	addressCoin, err := s.addressCoinRepository.Find(address.FindAddressCoinOptions{
		Network:    input.Network,
		CoinName:   input.CoinName,
		MnemonicId: input.MnemonicId,
	}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not find coin: %s", err))
	}
	if addressCoin == nil {
		return nil, app_error.InternalError(errors.Errorf("coin of network %s, not found", input.Network))
	}

	intAmount, err := app_util.AmountToInt(input.Network, input.Amount)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not convert amount to int: %s", err))
	}
	if addressCoin.Amount > intAmount {
		addressCoin.Amount = addressCoin.Amount - intAmount
		err := s.addressCoinRepository.Save(addressCoin, nil)
		if err != nil {
			return nil, app_error.InternalError(errors.Errorf("can not save coin address amount: %s", err))
		}
	}

	s.logger.Info(fmt.Sprintf("transfer of coin name %s and amount %d was made", input.CoinName, addressCoin.Amount))

	return &TransferOutput{Hash: result.Hash}, nil
}
