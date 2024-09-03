package coin

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type CoinRepository struct {
	baseRepository *repository.BaseRepository
}

func NewCoinRepository(baseRepository *repository.BaseRepository) *CoinRepository {
	return &CoinRepository{baseRepository: baseRepository}
}

func (r CoinRepository) FindAll(tx driver.Tx) ([]*Coin, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{},
		Limit:      50000,
		Offset:     0,
	}, &Coin{}, tx)
	if err != nil {
		return nil, fmt.Errorf("error find coins %s", err)
	}
	return items, nil
}

func (r CoinRepository) FindAllMappedByNetworks(tx driver.Tx) (map[app_enum.Network][]Coin, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{},
		Limit:      50000,
		Offset:     0,
	}, &Coin{}, tx)
	if err != nil {
		return nil, fmt.Errorf("error find coins %s", err)
	}
	coinsByNetwork := make(map[app_enum.Network][]Coin)
	for _, coinData := range items {
		coinsByNetwork[coinData.Network] = append(coinsByNetwork[coinData.Network], *coinData)
	}

	return coinsByNetwork, err
}

func (r CoinRepository) FindAllMappedByIds(tx driver.Tx) (map[int64]Coin, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{},
		Limit:      50000,
		Offset:     0,
	}, &Coin{}, tx)
	if err != nil {
		return nil, fmt.Errorf("error find coins %s", err)
	}
	coinsById := make(map[int64]Coin)
	for _, coinData := range items {
		coinsById[coinData.Id] = *coinData
	}

	return coinsById, nil
}
