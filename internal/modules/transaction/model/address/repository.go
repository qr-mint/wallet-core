package address

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	baseRepository *repository.BaseRepository
}

func NewRepository(baseRepository *repository.BaseRepository) *Repository {
	return &Repository{
		baseRepository,
	}
}

type FindManyAddressOptions struct {
	MnemonicId int64
}

func (r Repository) FindMany(options FindManyAddressOptions, tx driver.Tx) ([]*Address, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{"mnemonic_id": options.MnemonicId},
		Limit:      50000,
		Offset:     0,
		OrderBy:    nil,
	}, &Address{}, tx)
	if err != nil {
		return nil, fmt.Errorf("error while find many addresses: %w", err)
	}

	return items, nil
}
