package address

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	baseRepository *repository.BaseRepository
}

func NewRepository(baseRepository *repository.BaseRepository) *Repository {
	return &Repository{baseRepository}
}

type FindOptions struct {
	MnemonicId int64
	Network    app_enum.Network
}

func (r Repository) Find(options FindOptions, tx driver.Tx) (*Address, error) {
	address := Address{}
	err := r.baseRepository.FindOneBy(goqu.Ex{"mnemonic_id": options.MnemonicId, "network": options.Network}, &address, tx)
	if err != nil {
		return nil, err
	}

	return &address, nil
}

type FindManyOptions struct {
	MnemonicId int64
}

func (r Repository) FindMany(options FindManyOptions, tx driver.Tx) (map[int64]Address, error) {
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
