package coin

import (
	"database/sql/driver"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	baseRepository *repository.BaseRepository
}

func NewRepository(baseRepository *repository.BaseRepository) *Repository {
	return &Repository{baseRepository: baseRepository}
}

type FindOptions struct {
	Network  app_enum.Network
	CoinName app_enum.CoinName
}

func (r Repository) Find(options FindOptions, tx driver.Tx) (*Coin, error) {
	coin := Coin{}
	err := r.baseRepository.FindOneBy(goqu.Ex{"network": options.Network, "name": options.CoinName}, &coin, tx)
	if err != nil {
		return nil, errors.Errorf("can not find the coin by options %v: %s", options, err)
	}

	if coin.Id == 0 {
		return nil, nil
	}

	return &coin, err
}
