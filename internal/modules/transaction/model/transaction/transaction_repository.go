package transaction

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/internal/modules/transaction/enum"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	baseRepository *repository.BaseRepository
}

func NewTransactionRepository(baseRepository *repository.BaseRepository) *Repository {
	return &Repository{baseRepository: baseRepository}
}

type FindManyOptions struct {
	AddressCoinId *int64
	OnlyOut       *bool
	MnemonicId    int64
	Limit         uint
	Offset        uint
}

func (r Repository) FindMany(options FindManyOptions, tx driver.Tx) ([]*BlockchainTransaction, error) {
	expression := goqu.Ex{"address": goqu.Select("address").From("wallet_addresses").Where(goqu.Ex{"mnemonic_id": options.MnemonicId})}
	if options.AddressCoinId != nil {
		expression["coin_id"] = goqu.Select("coin_id").From("wallet_address_coins").Where(goqu.Ex{"id": options.AddressCoinId})
	}
	if options.OnlyOut != nil {
		expression["type"] = enum.OutType
	}

	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: expression,
		Limit:      options.Limit,
		Offset:     options.Offset,
		OrderBy:    goqu.I("created_at").Desc(),
	}, &BlockchainTransaction{}, tx)
	if err != nil {
		return nil, fmt.Errorf("error find transactions: %s", err)
	}
	return items, nil
}

func (r Repository) Save(transaction *BlockchainTransaction, tx driver.Tx) error {
	err := r.baseRepository.CreateOrUpdate(transaction, tx)
	if err != nil {
		return fmt.Errorf("can not create transaction: %s", err)
	}

	return nil
}
