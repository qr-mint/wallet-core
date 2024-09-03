package address

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type AddressCoinRepository struct {
	*repository.BaseRepository
}

func NewAddressCoinRepository(baseRepository *repository.BaseRepository) *AddressCoinRepository {
	return &AddressCoinRepository{
		baseRepository,
	}
}

type FindAddressCoinOptions struct {
	Network    app_enum.Network
	CoinName   app_enum.CoinName
	MnemonicId int64
}

func (r AddressCoinRepository) Find(options FindAddressCoinOptions, tx driver.Tx) (*AddressCoin, error) {
	addressCoin := AddressCoin{}
	err := r.FindOneBy(goqu.Ex{
		"coin_id":    goqu.Select("id").From("coins").Where(goqu.Ex{"network": options.Network, "name": options.CoinName}),
		"address_id": goqu.Select("id").From("wallet_addresses").Where(goqu.Ex{"mnemonic_id": options.MnemonicId}),
	}, &addressCoin, tx)
	if err != nil {
		return nil, fmt.Errorf("error while find one address coin: %w", err)
	}
	if addressCoin.Id == 0 {
		return nil, nil
	}

	return &addressCoin, nil
}

func (r AddressCoinRepository) Save(addressCoin *AddressCoin, tx driver.Tx) error {
	err := r.CreateOrUpdate(addressCoin, tx)
	if err != nil {
		return fmt.Errorf("can not save addressCoin: %s", err)
	}

	return nil
}
