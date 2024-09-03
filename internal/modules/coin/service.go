package coin

import (
	"fmt"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/internal/modules/coin/model/coin"
	"nexus-wallet/internal/modules/coin/model/price"
	"nexus-wallet/internal/modules/coin/provider"
)

type Service struct {
	priceRepository *price.Repository
	coinRepository  *coin.Repository
	syncer          *provider.Syncer
}

func NewService(
	repository *price.Repository,
	coinRepository *coin.Repository,
	syncer *provider.Syncer,
) *Service {
	service := &Service{
		coinRepository:  coinRepository,
		priceRepository: repository,
		syncer:          syncer,
	}

	return service
}

func (s Service) SyncPrices() {
	s.syncer.SyncAllPrices()
}

func (s Service) GetCryptoPriceInFiat(input GetCryptoPriceInFiatInput) (*GetCryptoPriceInFiatOutput, *app_error.AppError) {
	storedCoin, err := s.coinRepository.Find(coin.FindOptions{Network: input.Network, Name: input.CoinName}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get coin: %s", err))
	}
	if storedCoin == nil {
		return nil, app_error.InvalidDataError(fmt.Errorf("can not find coin %s in network %s", input.CoinName, input.Network))
	}
	storedPrice, err := s.priceRepository.FindLatest(price.FindOptions{
		CoinId:       storedCoin.Id,
		FiatCurrency: input.FiatCurrency,
	}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get rate from database: %s", err))
	}
	if storedPrice == nil {
		return nil, nil
	}

	floatPrice, err := app_util.AmountToFloat(input.Network, storedPrice.Price)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not transform price to int: %s", err))
	}

	return &GetCryptoPriceInFiatOutput{PayableAmount: floatPrice * input.Amount, Price: floatPrice}, nil
}

func (s Service) GetFiatPriceInCrypto(input GetFiatPriceInCryptoInput) (*GetFiatPriceInCryptoOutput, *app_error.AppError) {
	storedCoin, err := s.coinRepository.Find(coin.FindOptions{Network: input.Network, Name: input.CoinName}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get coin: %s", err))
	}
	if storedCoin == nil {
		return nil, app_error.InvalidDataError(fmt.Errorf("can not find coin %s in network %s", input.CoinName, input.Network))
	}
	storedPrice, err := s.priceRepository.FindLatest(price.FindOptions{
		CoinId:       storedCoin.Id,
		FiatCurrency: input.FiatCurrency,
	}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get rate from database: %s", err))
	}
	if storedPrice == nil {
		return nil, nil
	}

	floatPrice, err := app_util.AmountToFloat(input.Network, storedPrice.Price)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not transform price to int: %s", err))
	}

	return &GetFiatPriceInCryptoOutput{PayableAmount: input.Amount / floatPrice, Price: 1 / floatPrice}, nil
}
