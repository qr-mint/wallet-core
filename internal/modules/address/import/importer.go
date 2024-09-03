package _import

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/internal/modules/address/model/address"
	"nexus-wallet/internal/modules/address/model/coin"
	"nexus-wallet/internal/modules/address/model/mnemonic"
	"nexus-wallet/pkg/repository"
	"nexus-wallet/pkg/transaction"
)

type Importer struct {
	addressRepository     *address.AddressRepository
	addressCoinRepository *address.AddressCoinRepository
	mnemonicRepository    *mnemonic.Repository
	coinRepository        *coin.CoinRepository
	transactionManager    transaction.Manager
}

func NewImporter(
	addressRepository *address.AddressRepository,
	addressCoinRepository *address.AddressCoinRepository,
	mnemonicRepository *mnemonic.Repository,
	coinRepository *coin.CoinRepository,
	transactionManager transaction.Manager,
) *Importer {
	return &Importer{
		addressRepository:     addressRepository,
		addressCoinRepository: addressCoinRepository,
		mnemonicRepository:    mnemonicRepository,
		coinRepository:        coinRepository,
		transactionManager:    transactionManager,
	}
}

type AddressData struct {
	Network app_enum.Network
	Address string
}

type ImportData struct {
	Addresses    []AddressData
	Name         string
	MnemonicHash string
}

func (f Importer) Import(userId int64, data ImportData) *app_error.AppError {
	if data.MnemonicHash == "" {
		return app_error.InvalidDataError(errors.New("mnemonic hash must be provided"))
	}
	var appError *app_error.AppError
	err := f.transactionManager.WithTransaction(func(tx driver.Tx) error {
		mnemonicData, err := f.findExistingWalletMnemonicByHash(data, tx)
		if err != nil {
			return fmt.Errorf("could not find mnemonic: %w", err)
		}
		if mnemonicData != nil {
			err := f.assignWalletExistingWalletToUser(userId, mnemonicData, tx)
			if err != nil {
				return fmt.Errorf("could not assign wallet to user: %w", err)
			}
			return nil
		}
		appError = f.createNewWallet(userId, data, tx)
		if appError != nil {
			return appError.Error
		}

		return nil
	})
	if appError != nil {
		return appError
	}
	if err != nil {
		return app_error.InternalError(fmt.Errorf("error in transaction: %s", err))
	}

	return nil
}

func (f Importer) createNewWallet(userId int64, walletData ImportData, tx driver.Tx) *app_error.AppError {
	if walletData.Name == "" {
		return app_error.InvalidDataError(errors.New("wallet name must be provided"))
	}

	mnemonicModel := mnemonic.Mnemonic{Hash: walletData.MnemonicHash, Name: walletData.Name}
	err := f.mnemonicRepository.Save(&mnemonicModel, userId, tx)
	if err != nil {
		return app_error.InternalError(fmt.Errorf("can not save mnemonic while creating new wallet: %s", err))
	}

	coinsByNetwork, err := f.coinRepository.FindAllMappedByNetworks(tx)
	if err != nil {
		return app_error.InternalError(err)
	}

	createdAddressesNetworks := make(map[app_enum.Network]struct{})
	for _, addressData := range walletData.Addresses {
		if _, addressAlreadyCreated := createdAddressesNetworks[addressData.Network]; addressAlreadyCreated {
			return app_error.InvalidDataError(errors.New("invalid addresses data provided: multiple addresses in one network"))
		}

		if err := f.createAddress(mnemonicModel, addressData, coinsByNetwork, tx); err != nil {
			return err
		}

		createdAddressesNetworks[addressData.Network] = struct{}{}
	}
	if len(createdAddressesNetworks) != len(walletData.Addresses) {
		return app_error.InvalidDataError(errors.New("invalid addresses data provided: not all networks provided"))
	}

	return nil
}

func (f Importer) createAddress(mnemonicModel mnemonic.Mnemonic, addressData AddressData, coins map[app_enum.Network][]coin.Coin, tx driver.Tx) *app_error.AppError {
	if _, networkKeyExists := coins[addressData.Network]; !networkKeyExists {
		return app_error.InvalidDataError(fmt.Errorf("invalid network provided: %s", addressData.Network))
	}

	isValidAddress, err := app_util.IsValidAddress(addressData.Address, addressData.Network)
	if err != nil {
		return app_error.InvalidDataError(fmt.Errorf("can not validate address: %s", err))
	}
	if !isValidAddress {
		return app_error.InvalidDataError(fmt.Errorf("invalid address provided for %s network: %s", addressData.Network, addressData.Address))
	}

	addressModel := address.Address{Address: addressData.Address, Network: addressData.Network, MnemonicId: mnemonicModel.Id}
	appErr := f.addressRepository.Save(&addressModel, tx)
	if appErr != nil {
		if repository.IsUniqueError(appErr) {
			return app_error.InvalidDataError(
				fmt.Errorf("address of network %s is already exists with another mnemonic hash", addressData.Network),
			)
		}
		return app_error.InternalError(fmt.Errorf("can not create address: %s", appErr))
	}

	for _, coinData := range coins[addressData.Network] {
		addressCoin := address.AddressCoin{Amount: 0, Address: "", IsVisible: coinData.IsDefault, CoinId: coinData.Id, AddressId: addressModel.Id}
		appErr := f.addressCoinRepository.Save(&addressCoin, tx)
		if appErr != nil {
			return app_error.InternalError(fmt.Errorf("can not save address coin with coin %d: %s", coinData.Id, appErr))
		}
	}

	return nil
}

func (f Importer) assignWalletExistingWalletToUser(userId int64, mnemonicData *mnemonic.Mnemonic, tx driver.Tx) error {
	isAlreadyAssignedToUser, err := f.mnemonicRepository.IsAssignedToUser(mnemonic.IsAssignedToUserOptions{UserId: userId, MnemonicId: mnemonicData.Id}, tx)
	if err != nil {
		return fmt.Errorf("can not check mnemonic assigned to user: %s", err)
	}
	if isAlreadyAssignedToUser {
		return nil
	}
	err = f.mnemonicRepository.AssignUser(mnemonicData, userId, tx)
	if err != nil {
		return fmt.Errorf("assign user to existing mnemonic: %s", err)
	}

	return nil
}

func (f Importer) findExistingWalletMnemonicByHash(walletData ImportData, tx driver.Tx) (*mnemonic.Mnemonic, error) {
	mnemonicData, err := f.mnemonicRepository.Find(mnemonic.FindOptions{Hash: walletData.MnemonicHash}, tx)
	if err != nil {
		return nil, fmt.Errorf("can not find mnemonic in assign user check: %s", err)
	}
	if mnemonicData == nil {
		return nil, nil
	}

	return mnemonicData, nil
}
