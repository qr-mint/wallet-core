package exchange

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
	return &Repository{baseRepository}
}

type FindManyOptions struct {
	MnemonicId int64
	Limit      uint
	Offset     uint
}

func (r Repository) FindMany(options FindManyOptions, tx driver.Tx) ([]*Exchange, error) {
	models, err := repository.FindManyBy(
		r.baseRepository,
		repository.FindManyByOptions{
			Expression: goqu.Ex{"mnemonic_id": options.MnemonicId},
			Limit:      options.Limit,
			Offset:     options.Offset,
		},
		&Exchange{},
		tx,
	)

	return models, err
}

func (r Repository) Save(exchange *Exchange, tx driver.Tx) error {
	err := r.baseRepository.CreateOrUpdate(exchange, tx)
	if err != nil {
		return fmt.Errorf("can not save exchange: %s", err)
	}

	return nil
}
