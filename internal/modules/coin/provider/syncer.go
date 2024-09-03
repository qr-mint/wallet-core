package provider

import (
	"fmt"
	"gitlab.com/golib4/logger/logger"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/internal/modules/coin/model/coin"
	"nexus-wallet/internal/modules/coin/model/price"
	"time"
)

type Params struct {
	FrequencyInHours int
}

type Syncer struct {
	priceRepository *price.Repository
	coinRepository  *coin.Repository
	logger          logger.Logger
	providers       []Provider
	params          Params
}

func NewSyncer(
	priceRepository *price.Repository,
	coinRepository *coin.Repository,
	logger logger.Logger,
	providers []Provider,
	params Params,
) *Syncer {
	return &Syncer{
		priceRepository: priceRepository,
		coinRepository:  coinRepository,
		logger:          logger,
		providers:       providers,
		params:          params,
	}
}

func (s Syncer) SyncAllPrices() {
	for _, currency := range app_enum.GetCurrencies() {
		for _, network := range app_enum.GetNetworks() {
			for _, coinName := range app_enum.GetCoinNamesByNetwork(network) {
				err := s.syncPriceFromProvider(network, coinName, currency, 0)
				if err != nil {
					s.logger.Warningf("Error syncing price: %s", err)
				}

				time.Sleep(10 * time.Second)
			}
		}
	}
}

func (s Syncer) syncPriceFromProvider(network app_enum.Network, name app_enum.CoinName, currency app_enum.Currency, attempts int) error {
	priceProvider, err := s.getProvider(network, name)
	if err != nil {
		return fmt.Errorf("can not get provider: %s", err)
	}
	priceAmount, err := priceProvider.GetPrice(currency, time.Now())
	if err != nil {
		if attempts <= 2 {
			time.Sleep(10 * time.Second)
			return s.syncPriceFromProvider(network, name, currency, attempts+1)
		}

		return fmt.Errorf("can not get price from provider: %s", err)
	}
	intPrice, err := app_util.AmountToInt(network, priceAmount)
	if err != nil {
		return fmt.Errorf("can not transform price to int: %s", err)
	}
	storedCoin, err := s.coinRepository.Find(coin.FindOptions{Network: network, Name: name}, nil)
	if err != nil {
		return fmt.Errorf("can not get coin: %s", err)
	}
	if storedCoin == nil {
		return fmt.Errorf("can not find coin %s in network %s", name, network)
	}
	err = s.priceRepository.Save(&price.Price{
		Price:        intPrice,
		CoinId:       storedCoin.Id,
		FiatCurrency: string(currency),
		Date:         time.Now(),
	}, nil)
	if err != nil {
		return fmt.Errorf("can not save price: %s", err)
	}

	return nil
}

func (s Syncer) getProvider(network app_enum.Network, name app_enum.CoinName) (Provider, error) {
	for _, priceProvider := range s.providers {
		if priceProvider.IsSupports(network, name) {
			return priceProvider, nil
		}
	}

	return nil, fmt.Errorf("network %s and coin %s not supports", network, name)
}
