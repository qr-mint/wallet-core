package coin

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	*repository.BaseRepository
}

func NewRepository(repository *repository.BaseRepository) *Repository {
	return &Repository{repository}
}

type FindOptions struct {
	Network app_enum.Network
	Name    app_enum.CoinName
}

func (r Repository) Find(options FindOptions, tx driver.Tx) (*Coin, error) {
	coin := Coin{}
	err := r.FindOneBy(goqu.Ex{"network": options.Network, "name": options.Name}, &coin, tx)
	if err != nil {
		return nil, fmt.Errorf("can not find coin by options %v: %w", options, err)
	}
	if coin.Id == 0 {
		return nil, nil
	}

	return &coin, nil
}
