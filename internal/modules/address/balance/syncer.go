package balance

import (
	"fmt"
	"gitlab.com/golib4/logger/logger"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/address/balance/provider"
	"nexus-wallet/internal/modules/address/model/address"
	"nexus-wallet/internal/modules/address/model/coin"
	sync_manager "nexus-wallet/internal/shared/sync"
	"sync"
)

type Params struct {
	MaxAttempts int32
}

type Syncer struct {
	addressCoinRepository *address.AddressCoinRepository
	addressRepository     *address.AddressRepository
	coinRepository        *coin.CoinRepository
	logger                logger.Logger
	providers             map[app_enum.Network][]provider.Provider
	params                Params
	syncManager           *sync_manager.Manager
}

func NewSyncer(
	addressCoinRepository *address.AddressCoinRepository,
	addressRepository *address.AddressRepository,
	coinRepository *coin.CoinRepository,
	logger logger.Logger,
	providers map[app_enum.Network][]provider.Provider,
	params Params,
	syncManager *sync_manager.Manager,
) *Syncer {
	return &Syncer{
		addressCoinRepository: addressCoinRepository,
		addressRepository:     addressRepository,
		coinRepository:        coinRepository,
		logger:                logger,
		providers:             providers,
		params:                params,
		syncManager:           syncManager,
	}
}

const baseSyncKey = "balances_last_sync_mnemonic_"

func (s Syncer) Sync(mnemonicId int64) error {
	lastSyncTime, err := s.syncManager.GetLastSyncTime(baseSyncKey, mnemonicId)
	if err != nil {
		return fmt.Errorf("can not get timestamp for sync %s", err)
	}
	if !s.syncManager.NeedToSync(lastSyncTime) {
		return nil
	}

	addresses, err := s.addressRepository.FindMany(address.FindManyAddressOptions{
		MnemonicId: mnemonicId,
	}, nil)
	if err != nil {
		return fmt.Errorf("can not get addresses %w", err)
	}
	var wg sync.WaitGroup
	wg.Add(len(addresses))

	allDataSynced := true
	for _, addressData := range addresses {
		go func(a address.Address) {
			defer wg.Done()
			err := s.syncBalanceByAddress(a, 0)
			if err != nil {
				allDataSynced = false
				s.logger.Warningf("can not sync balance %s", err.Error)
				return
			}
		}(addressData)
	}
	wg.Wait()

	if allDataSynced {
		err = s.syncManager.SetLastSyncTime(baseSyncKey, mnemonicId)
		if err != nil {
			return fmt.Errorf("can not set last sync time %s", err)
		}
	}

	return nil
}

func (s Syncer) syncBalanceByAddress(addressData address.Address, attempts int32) *app_error.AppError {
	addressCoinsData, err := s.addressCoinRepository.FindManyByAddress(address.FindManyByAddressOptions{AddressId: addressData.Id}, nil)
	if err != nil {
		return app_error.InternalError(fmt.Errorf("can not get wallet address %w", err))
	}
	coinsData, err := s.coinRepository.FindAllMappedByIds(nil)
	if err != nil {
		return app_error.InternalError(fmt.Errorf("can not get coins mapped by ids in syncBalanceByNetwork %s", err))
	}
	balanceProvider, err := s.getProvider(addressData.Network)
	if err != nil {
		return app_error.InternalError(fmt.Errorf("can not get balance provider %w", err))
	}
	for _, addressCoinData := range addressCoinsData {
		coinData := coinsData[addressCoinData.CoinId]
		var balanceData *provider.GetBalanceOutput
		var appErr *app_error.AppError
		if coinData.IsToken {
			balanceData, appErr = balanceProvider.GetCoinBalance(addressData.Address, *coinData.Address)
		} else {
			balanceData, appErr = balanceProvider.GetMainCoinBalance(addressData.Address)
		}
		if appErr != nil {
			if attempts > s.params.MaxAttempts || appErr.Code != app_error.Internal {
				return appErr
			}

			s.logger.Warningf("can not get balance from provider attempt %s", appErr.Error)
			return s.syncBalanceByAddress(addressData, attempts+1)
		}
		addressCoinData.Amount = balanceData.Value

		err = s.addressCoinRepository.Save(addressCoinData, nil)
		if err != nil {
			return app_error.InternalError(fmt.Errorf("can not save address coin amount: %w", err))
		}
	}
	return nil
}

var lastUsedProviderIndex = 0

func (s Syncer) getProvider(network app_enum.Network) (provider.Provider, error) {
	balanceProviders, exists := s.providers[network]
	if !exists {
		return nil, fmt.Errorf("provider %s is not supported", network)
	}

	index := lastUsedProviderIndex + 1
	if index >= len(balanceProviders) {
		index = 0
	}
	lastUsedProviderIndex = index

	return balanceProviders[index], nil
}
