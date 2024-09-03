package provider

import (
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/golib4/logger/logger"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/nft/model/address"
	"nexus-wallet/internal/modules/nft/model/nft"
	sync_manager "nexus-wallet/internal/shared/sync"
	"nexus-wallet/pkg/transaction"
)

type Params struct {
	MaxAttempts int32
}

type Syncer struct {
	syncManager        *sync_manager.Manager
	logger             logger.Logger
	addressRepository  *address.Repository
	nftRepository      *nft.Repository
	nftProviders       map[app_enum.Network]NftProvider
	transactionManager transaction.Manager
	params             Params
}

func NewSyncer(
	syncManager *sync_manager.Manager,
	logger logger.Logger,
	addressRepository *address.Repository,
	nftRepository *nft.Repository,
	nftProviders map[app_enum.Network]NftProvider,
	transactionManager transaction.Manager,
	params Params,
) *Syncer {
	return &Syncer{
		syncManager:        syncManager,
		logger:             logger,
		addressRepository:  addressRepository,
		nftRepository:      nftRepository,
		nftProviders:       nftProviders,
		transactionManager: transactionManager,
		params:             params,
	}
}

const baseSyncKey = "nft_last_sync_mnemonic_"

func (s Syncer) Sync(mnemonicId int64) error {
	lastSyncTime, err := s.syncManager.GetLastSyncTime(baseSyncKey, mnemonicId)
	if err != nil {
		return fmt.Errorf("can not get timestamp for sync %s", err)
	}
	if !s.syncManager.NeedToSync(lastSyncTime) {
		return nil
	}

	addressData, err := s.addressRepository.Find(address.FindOptions{MnemonicId: mnemonicId, Network: app_enum.TonNetwork}, nil)
	if err != nil {
		return fmt.Errorf("can not get addresses %w", err)
	}
	syncErr := s.syncNftByAddress(mnemonicId, *addressData, 0)
	if syncErr != nil {
		s.logger.Warningf("can not sync nft %s", syncErr.Error)
		return nil
	}

	err = s.syncManager.SetLastSyncTime(baseSyncKey, mnemonicId)
	if err != nil {
		return fmt.Errorf("can not set last sync time %s", err)
	}

	return nil
}

func (s Syncer) syncNftByAddress(mnemonicId int64, addressData address.Address, attempts int32) *app_error.AppError {
	nftProvider, nftProviderExists := s.nftProviders[addressData.Network]
	if !nftProviderExists {
		return app_error.InternalError(errors.Errorf("nft provider for network %s does not exists", addressData.Network))
	}
	data, provideErr := nftProvider.Provide(ProvideInput{OwnerAddress: addressData.Address})
	if provideErr != nil {
		if attempts > s.params.MaxAttempts || provideErr.Code != app_error.Internal {
			return provideErr
		}

		return s.syncNftByAddress(mnemonicId, addressData, attempts+1)
	}

	err := s.transactionManager.WithTransaction(func(tx driver.Tx) error {
		err := s.nftRepository.DeleteBy(nft.DeleteOptions{AddressId: addressData.Id}, tx)
		if err != nil {
			return errors.Errorf("can not delete nfts %s", err)
		}
		for _, nftItem := range data.Items {
			err := s.nftRepository.Save(&nft.Nft{
				Address:               nftItem.Address,
				Name:                  nftItem.Name,
				Price:                 nftItem.Price,
				TokenSymbol:           nftItem.PriceTokenSymbol,
				Index:                 nftItem.Index,
				CollectionAddress:     nftItem.CollectionAddress,
				CollectionName:        nftItem.CollectionName,
				CollectionDescription: nftItem.CollectionDescription,
				PreviewData:           nftItem.PreviewsData,
				AddressId:             addressData.Id,
			}, tx)
			if err != nil {
				return errors.Errorf("can not save nft %s", err)
			}
		}

		return nil
	})
	if err != nil {
		return app_error.InternalError(errors.Errorf("can not sync nft by address (%s): %s", addressData.Address, err))
	}

	return nil
}
