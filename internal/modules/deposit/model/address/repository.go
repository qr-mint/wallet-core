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
	query, _, err := goqu.Select(
		goqu.Select("address").From("wallet_addresses").Where(
			goqu.L("wallet_addresses.id = wallet_address_coins.address_id"),
			goqu.Ex{"wallet_addresses.mnemonic_id": options.MnemonicId},
		),
		goqu.Select("name").From("coins").Where(
			goqu.L("coins.id = wallet_address_coins.coin_id"),
		),
	).From("wallet_address_coins").Where(goqu.Ex{"id": options.Id}).ToSQL()
	if err != nil {
		return nil, errors.Errorf("generate address coin %d query error: %s", options.Id, err)
	}
	row, err := r.QueryRow(tx, query)
	if err != nil {
		return nil, errors.Errorf("query address coin %d error: %s", options.Id, err)
	}

	var address *string
	var coinName *app_enum.CoinName
	err = row.Scan(&address, &coinName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("can not scan address coin query result: %s", err)
	}

	if address == nil {
		return nil, nil
	}

	return &AddressCoin{
		WalletAddress: *address,
		CoinName:      *coinName,
	}, nil
}
