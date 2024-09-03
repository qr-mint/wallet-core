package provider

import (
	"fmt"
	"gitlab.com/golib4/logger/logger"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/transaction/model/address"
	"nexus-wallet/internal/modules/transaction/model/coin"
	"nexus-wallet/internal/modules/transaction/model/transaction"
	sync_manager "nexus-wallet/internal/shared/sync"
	"nexus-wallet/pkg/repository"
	"sync"
	"time"
)

type Syncer struct {
	transactionRepository *transaction.Repository
	addressRepository     *address.Repository
	coinRepository        *coin.Repository
	logger                logger.Logger
	providers             map[app_enum.Network]Provider
	syncManager           *sync_manager.Manager
}

func NewSyncer(
	transactionRepository *transaction.Repository,
	addressRepository *address.Repository,
	coinRepository *coin.Repository,
	logger logger.Logger,
	providers map[app_enum.Network]Provider,
	syncManager *sync_manager.Manager,
) *Syncer {
	return &Syncer{
		transactionRepository: transactionRepository,
		addressRepository:     addressRepository,
		coinRepository:        coinRepository,
		logger:                logger,
		providers:             providers,
		syncManager:           syncManager,
	}
}

const baseSyncKey = "transactions_last_sync_mnemonic_"

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
		return fmt.Errorf("can not get addresses  %s", err)
	}

	coins, err := s.coinRepository.FindAllMappedByNames(nil)
	if err != nil {
		return fmt.Errorf("can not get coins %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(addresses))

	allDataSynced := true
	for _, addressData := range addresses {
		go func(a address.Address) {
			defer wg.Done()
			err := s.syncTransactionsByAddress(a, coins, lastSyncTime.Timestamp.Add(-1*time.Hour))
			if err != nil {
				s.logger.Warningf("can not sync transactions %s", err)
				allDataSynced = false
				return
			}
		}(*addressData)
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

func (s Syncer) syncTransactionsByAddress(
	addressData address.Address,
	coins map[app_enum.CoinName]coin.Coin,
	timestampFrom time.Time,
) error {
	transactionProvider, providerExists := s.providers[addressData.Network]
	if !providerExists {
		return fmt.Errorf("can not get transaction provider by network: %s", addressData.Network)
	}
	transactions, err := transactionProvider.GetTransactions(addressData.Address, 100, timestampFrom)
	if err != nil {
		return fmt.Errorf("can not get transactions %s", err)
	}
	for _, transactionData := range transactions {
		err := s.transactionRepository.Save(&transaction.BlockchainTransaction{
			Hash:        transactionData.Hash,
			Amount:      transactionData.Amount,
			Address:     addressData.Address,
			AddressTo:   transactionData.To,
			AddressFrom: transactionData.From,
			Status:      transactionData.Status,
			Type:        transactionData.Type,
			CoinId:      coins[transactionData.CoinName].Id,
			CreatedAt:   transactionData.CreatedAt,
		}, nil)
		if err != nil {
			if repository.IsUniqueError(err) {
				continue
			}

			return fmt.Errorf("can not save transaction: %w", err)
		}
	}

	return nil
}
