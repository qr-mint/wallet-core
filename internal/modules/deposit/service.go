package deposit

import (
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/deposit/enum"
	"nexus-wallet/internal/modules/deposit/model/address"
	"nexus-wallet/internal/modules/deposit/provider"
)

type Service struct {
	providers         map[enum.ProviderName]provider.RedirectProvider
	addressRepository *address.Repository
}

func NewService(providers map[enum.ProviderName]provider.RedirectProvider, addressRepository *address.Repository) *Service {
	return &Service{
		providers:         providers,
		addressRepository: addressRepository,
	}
}

func (s Service) GetLimits(input GetLimitsInput) (*GetLimitsOutput, *app_error.AppError) {
	addressData, err := s.addressRepository.Find(address.FindOptions{Id: input.AddressCoinId, MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("cannot find address: %s", err))
	}
	if addressData == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("address coin with id %d not found", input.AddressCoinId))
	}

	redirectProvider, providerExists := s.providers[input.ProviderName]
	if !providerExists {
		return nil, app_error.InternalError(errors.Errorf("failed to get redirect provider for limits: %s", redirectProvider))
	}
	limits, providerErr := redirectProvider.GetLimits(addressData.CoinName, input.FiatCurrency)
	if providerErr != nil {
		return nil, providerErr
	}

	return &GetLimitsOutput{Min: limits.Min, Max: limits.Max}, nil
}

func (s Service) ProvideRedirectUrl(input ProvideRedirectUrlInput) (*ProvideRedirectUrlOutput, *app_error.AppError) {
	addressData, err := s.addressRepository.Find(address.FindOptions{Id: input.AddressCoinId, MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("cannot find address: %s", err))
	}
	if addressData == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("address with id %d not found", input.AddressCoinId))
	}

	redirectProvider, providerExists := s.providers[input.ProviderName]
	if !providerExists {
		return nil, app_error.InternalError(errors.Errorf("failed to get redirect provider: %s", redirectProvider))
	}
	url, providerErr := redirectProvider.ProvideRedirectUrl(addressData.CoinName, input.FiatCurrency, input.Amount, addressData.WalletAddress)
	if providerErr != nil {
		return nil, providerErr
	}

	return &ProvideRedirectUrlOutput{url}, nil
}
