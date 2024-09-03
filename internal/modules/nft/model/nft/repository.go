package nft

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	baseRepository *repository.BaseRepository
}

func NewRepository(baseRepository *repository.BaseRepository) *Repository {
	return &Repository{
		baseRepository: baseRepository,
	}
}

type FindOptions struct {
	MnemonicId int64
	Id         int64
}

func (r Repository) Find(options FindOptions, tx driver.Tx) (*Nft, error) {
	nft := Nft{}
	err := r.baseRepository.FindOneBy(
		goqu.Ex{
			"id":         options.Id,
			"address_id": goqu.Select("id").From("wallet_addresses").Where(goqu.Ex{"mnemonic_id": options.MnemonicId}),
		}, &nft, tx)
	if err != nil {
		return nil, errors.Errorf("find nft(%d), error: %s", options.Id, err)
	}
	if nft.Id == 0 {
		return nil, nil
	}

	return &nft, nil
}

func (r Repository) FindOne(id int64, tx driver.Tx) (*Nft, error) {
	nft := Nft{}
	err := r.baseRepository.FindOne(id, &nft, tx)
	if err != nil {
		return nil, errors.Errorf("find one nft(%d), error: %s", id, err)
	}
	if nft.Id == 0 {
		return nil, nil
	}

	return &nft, nil
}

type FindManyOptions struct {
	MnemonicId int64
	Limit      uint
	Offset     uint
}

func (r Repository) FindMany(options FindManyOptions, tx driver.Tx) ([]*Nft, error) {
	models, err := repository.FindManyBy(
		r.baseRepository,
		repository.FindManyByOptions{
			Expression: goqu.Ex{"address_id": goqu.Select("id").From("wallet_addresses").Where(goqu.Ex{"mnemonic_id": options.MnemonicId})},
			Limit:      options.Limit,
			Offset:     options.Offset,
		},
		&Nft{},
		tx,
	)

	return models, err
}

func (r Repository) Save(nft *Nft, tx driver.Tx) error {
	err := r.baseRepository.CreateOrUpdate(nft, tx)
	if err != nil {
		return fmt.Errorf("can not save nft: %s", err)
	}

	return nil
}

type DeleteOptions struct {
	AddressId int64
}

func (r Repository) DeleteBy(options DeleteOptions, tx driver.Tx) error {
	err := r.baseRepository.DeleteBy(goqu.Ex{"address_id": options.AddressId}, &Nft{}, tx)
	if err != nil {
		return fmt.Errorf("can not delete nft by: %s", err)
	}

	return nil
}

func (r Repository) Delete(nft *Nft, tx driver.Tx) error {
	err := r.baseRepository.DeleteBy(goqu.Ex{"id": nft.Id}, &Nft{}, tx)
	if err != nil {
		return fmt.Errorf("can not delete nft: %s", err)
	}

	return nil
}
