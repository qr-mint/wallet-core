package price

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

func NewRepository(repository *repository.BaseRepository) *Repository {
	return &Repository{repository}
}

type FindOptions struct {
	CoinId       int64
	FiatCurrency app_enum.Currency
}

func (r Repository) FindLatest(options FindOptions, tx driver.Tx) (*Price, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{"coin_id": options.CoinId, "fiat_currency": options.FiatCurrency},
		Limit:      1,
		Offset:     0,
		OrderBy:    goqu.I("date").Desc(),
	}, &Price{}, tx)
	if err != nil {
		return nil, fmt.Errorf("error while finding coin price: %s", err)
	}

	if len(items) == 0 {
		return nil, nil
	}
	return items[0], nil
}

func (r Repository) Save(price *Price, tx driver.Tx) error {
	err := r.baseRepository.CreateOrUpdate(price, tx)
	if err != nil {
		return fmt.Errorf("can not save price: %s", err)
	}
	err = r.baseRepository.Refresh(price, tx)
	if err != nil {
		return fmt.Errorf("can not refresh price: %s", err)
	}

	return nil
}
