package address

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/pkg/repository"
)

type AddressRepository struct {
	baseRepository *repository.BaseRepository
}

func NewAddressRepository(baseRepository *repository.BaseRepository) *AddressRepository {
	return &AddressRepository{
		baseRepository,
	}
}

type FindManyAddressOptions struct {
	MnemonicId int64
}

func (r AddressRepository) FindMany(options FindManyAddressOptions, tx driver.Tx) (map[int64]Address, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{"mnemonic_id": options.MnemonicId},
		Limit:      50000,
		Offset:     0,
		OrderBy:    nil,
	}, &Address{}, tx)
	if err != nil {
		return nil, fmt.Errorf("error while find many addresses: %w", err)
	}

	addressesByIds := make(map[int64]Address)
	for _, addressData := range items {
		addressesByIds[addressData.Id] = *addressData
	}

	return addressesByIds, nil
}

func (r AddressRepository) Save(address *Address, tx driver.Tx) error {
	err := r.baseRepository.CreateOrUpdate(address, tx)
	if err != nil {
		return fmt.Errorf("can not create address: %s", err)
	}
	err = r.baseRepository.Refresh(address, tx)
	if err != nil {
		return fmt.Errorf("can not refresh address: %s", err)
	}

	return nil
}
