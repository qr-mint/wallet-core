package address

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	*repository.BaseRepository
}

func NewRepository(baseRepository *repository.BaseRepository) *Repository {
	return &Repository{baseRepository}
}

type FindOptions struct {
	Id         int64
	MnemonicId int64
}

func (r Repository) Find(options FindOptions, tx driver.Tx) (*AddressCoin, error) {
	query, _, err := goqu.
		From("wallet_address_coins").
		Select(
			"wallet_address_coins.id",
			"wallet_addresses.address",
			"coins.name",
			"coins.network",
		).
		Join(
			goqu.T("wallet_addresses"),
			goqu.On(goqu.Ex{
				"wallet_addresses.id": goqu.I("wallet_address_coins.address_id"),
			}),
		).
		Join(
			goqu.T("coins"),
			goqu.On(goqu.Ex{
				"coins.id": goqu.I("wallet_address_coins.coin_id"),
			}),
		).
		Where(goqu.Ex{
			"wallet_address_coins.id":      options.Id,
			"wallet_addresses.mnemonic_id": options.MnemonicId,
		}).
		ToSQL()
	if err != nil {
		return nil, errors.Errorf("generate address coin %d query error: %s", options.Id, err)
	}

	row, err := r.QueryRow(tx, query)
	if err != nil {
		return nil, errors.Errorf("query address coin %d error: %s", options.Id, err)
	}

	var addressCoinId *int64
	var address *string
	var coinName *app_enum.CoinName
	var coinNetwork *app_enum.Network
	err = row.Scan(&addressCoinId, &address, &coinName, &coinNetwork)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("can not scan address coin query result: %s", err)
	}

	if address == nil {
		return nil, nil
	}

	return &AddressCoin{
		AddressCoinId: *addressCoinId,
		WalletAddress: *address,
		CoinName:      *coinName,
		CoinNetwork:   *coinNetwork,
	}, nil
}
