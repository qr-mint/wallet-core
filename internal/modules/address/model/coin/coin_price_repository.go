package coin

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/pkg/repository"
	"time"
)

type PriceRepository struct {
	baseRepository *repository.BaseRepository
}

func NewPriceRepository(repository *repository.BaseRepository) *PriceRepository {
	return &PriceRepository{repository}
}

type FindOptions struct {
	CoinId       int64
	FiatCurrency app_enum.Currency
}

func (r PriceRepository) FindLatest(options FindOptions, tx driver.Tx) (*Price, error) {
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

func (r PriceRepository) FindYesterday(options FindOptions, tx driver.Tx) (*Price, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expressions: []goqu.Expression{
			goqu.Ex{
				"coin_id":       options.CoinId,
				"fiat_currency": options.FiatCurrency,
			},
			goqu.I("date").Lt(app_util.TimeToStartOfDay(time.Now())),
		},
		Limit:   1,
		Offset:  0,
		OrderBy: goqu.I("date").Desc(),
	}, &Price{}, tx)
	if err != nil {
		return nil, fmt.Errorf("error while finding coin price: %s", err)
	}

	if len(items) == 0 {
		return nil, nil
	}
	return items[0], nil
}
