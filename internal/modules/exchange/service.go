package exchange

import (
	"fmt"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/exchange/model/address"
	"nexus-wallet/internal/modules/exchange/model/coin"
	"nexus-wallet/internal/modules/exchange/model/exchange"
	"nexus-wallet/internal/modules/exchange/provider"
)

type Service struct {
	provider           provider.TransferableProvider
	addressRepository  *address.Repository
	exchangeRepository *exchange.Repository
	coinRepository     *coin.Repository
}

func NewService(
	provider provider.TransferableProvider,
	addressRepository *address.Repository,
	exchangeRepository *exchange.Repository,
	coinRepository *coin.Repository,
) *Service {
	return &Service{
		provider:           provider,
		addressRepository:  addressRepository,
		exchangeRepository: exchangeRepository,
		coinRepository:     coinRepository,
	}
}

func (s Service) List(input ListInput) (*ListOutput, *app_error.AppError) {
	items, err := s.exchangeRepository.FindMany(exchange.FindManyOptions{
		MnemonicId: input.MnemonicId,
		Limit:      input.Limit,
		Offset:     input.Offset,
	}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to get list exchange items: %s", err))
	}
	coinsMappedByIds, err := s.coinRepository.FindAllMappedByIds(nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get mapped coins by ids: %s", err))
	}

	return ListOutput{}.fillFromModel(items, coinsMappedByIds), nil
}

func (s Service) ProvideAddressForTransfer(input ProvideAddressForTransferInput) (*ProvideAddressForTransferOutput, *app_error.AppError) {
	addressCoinFrom, err := s.addressRepository.Find(address.FindOptions{Id: input.AddressCoinIdFrom, MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("cannot find from address: %s", err))
	}
	if addressCoinFrom == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("address coin from with id %d not found", input.AddressCoinIdFrom))
	}

	addressCoinTo, err := s.addressRepository.Find(address.FindOptions{Id: input.AddressCoinIdTo, MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("cannot find to address: %s", err))
	}
	if addressCoinTo == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("address coin to with id %d not found", input.AddressCoinIdTo))
	}

	data, provideErr := s.provider.ProvideAddressForTransfer(provider.ProvideAddressForTransferInput{
		CoinFromName:    addressCoinFrom.CoinName,
		CoinFromNetwork: addressCoinFrom.CoinNetwork,
		CoinToName:      addressCoinTo.CoinName,
		CoinToNetwork:   addressCoinTo.CoinNetwork,
		AddressTo:       addressCoinTo.WalletAddress,
		Amount:          input.Amount,
	})
	if provideErr != nil {
		return nil, provideErr
	}

	err = s.exchangeRepository.Save(&exchange.Exchange{
		ExternalId:  data.TransactionId,
		SupportLink: data.SupportLink,
		CoinFromId:  addressCoinFrom.AddressCoinId,
		CoinToId:    addressCoinTo.AddressCoinId,
		MnemonicId:  input.MnemonicId,
	}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("cannot save exchange data: %s", err))
	}

	return &ProvideAddressForTransferOutput{PayInAddress: data.PayInAddress, TransactionId: data.TransactionId}, nil
}

func (s Service) GetExchangeAmount(input GetExchangeAmountInput) (*GetExchangeAmountOutput, *app_error.AppError) {
	addressDataFrom, err := s.addressRepository.Find(address.FindOptions{Id: input.AddressCoinIdFrom, MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("cannot find from address: %s", err))
	}
	if addressDataFrom == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("address coin from with id %d not found", input.AddressCoinIdFrom))
	}

	addressDataTo, err := s.addressRepository.Find(address.FindOptions{Id: input.AddressCoinIdTo, MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("cannot find to address: %s", err))
	}
	if addressDataTo == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("address coin to with id %d not found", input.AddressCoinIdTo))
	}

	data, provideErr := s.provider.GetExchangeAmount(provider.GetExchangeAmountInput{
		CoinFromName:    addressDataFrom.CoinName,
		CoinFromNetwork: addressDataFrom.CoinNetwork,
		CoinToName:      addressDataTo.CoinName,
		CoinToNetwork:   addressDataTo.CoinNetwork,
		SendAmount:      input.SendAmount,
	})
	if provideErr != nil {
		return nil, provideErr
	}

	return &GetExchangeAmountOutput{ReceiveAmount: data.ReceiveAmount}, nil
}

func (s Service) GetLimits(input GetLimitsInput) (*GetLimitsOutput, *app_error.AppError) {
	addressDataFrom, err := s.addressRepository.Find(address.FindOptions{Id: input.AddressCoinIdFrom, MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("cannot find from address: %s", err))
	}
	if addressDataFrom == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("address coin from with id %d not found", input.AddressCoinIdFrom))
	}

	addressDataTo, err := s.addressRepository.Find(address.FindOptions{Id: input.AddressCoinIdTo, MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("cannot find to address: %s", err))
	}
	if addressDataTo == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("address coin to with id %d not found", input.AddressCoinIdFrom))
	}

	data, provideErr := s.provider.GetLimits(provider.GetLimitsInput{
		CoinFromName:    addressDataFrom.CoinName,
		CoinFromNetwork: addressDataFrom.CoinNetwork,
		CoinToName:      addressDataTo.CoinName,
		CoinToNetwork:   addressDataTo.CoinNetwork,
	})
	if provideErr != nil {
		return nil, provideErr
	}

	return &GetLimitsOutput{Min: data.Min, Max: data.Max}, nil
}
