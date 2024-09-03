package address

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/pkg/repository"
)

type AddressCoinRepository struct {
	baseRepository *repository.BaseRepository
}

func NewAddressCoinRepository(baseRepository *repository.BaseRepository) *AddressCoinRepository {
	return &AddressCoinRepository{
		baseRepository,
	}
}

type FindManyByAddressOptions struct {
	AddressId int64
}

func (r AddressCoinRepository) FindManyByAddress(options FindManyByAddressOptions, tx driver.Tx) ([]*AddressCoin, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{"address_id": options.AddressId},
		Limit:      50000,
		Offset:     0,
	}, &AddressCoin{}, tx)
	if err != nil {
		return nil, fmt.Errorf("find many by address coin error: %s", err)
	}

	return items, nil

}

type FindManyOptions struct {
	MnemonicId  int64
	OnlyVisible bool
}

func (r AddressCoinRepository) FindMany(options FindManyOptions, tx driver.Tx) ([]*AddressCoin, error) {
	expression := goqu.Ex{
		"address_id": goqu.Select("id").From("wallet_addresses").Where(goqu.Ex{"mnemonic_id": options.MnemonicId}),
	}
	if options.OnlyVisible {
		expression["is_visible"] = true
	}
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: expression,
		Limit:      50000,
		Offset:     0,
		OrderBy:    goqu.I("id").Desc(),
	}, &AddressCoin{}, tx)
	if err != nil {
		return nil, fmt.Errorf("find many by address coin error: %s", err)
	}

	return items, nil
}

type FindOptions struct {
	Id         int64
	MnemonicId int64
}

func (r AddressCoinRepository) Find(options FindOptions, tx driver.Tx) (*AddressCoin, error) {
	addressCoin := AddressCoin{}
	err := r.baseRepository.FindOneBy(goqu.Ex{
		"id":         options.Id,
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
	err := r.baseRepository.CreateOrUpdate(addressCoin, tx)
	if err != nil {
		return fmt.Errorf("can not save addressCoin: %s", err)
	}
	err = r.baseRepository.Refresh(addressCoin, tx)
	if err != nil {
		return fmt.Errorf("can not refresh addressCoin: %s", err)
	}

	return nil
}
