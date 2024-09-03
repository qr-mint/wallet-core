package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	baseRepository *repository.BaseRepository
}

func NewRepository(repository *repository.BaseRepository) *Repository {
	return &Repository{baseRepository: repository}
}

type FindOptions struct {
	Hash string
}

func (r Repository) Find(options FindOptions, tx driver.Tx) (*Mnemonic, error) {
	mnemonic := Mnemonic{}
	err := r.baseRepository.FindOneBy(goqu.Ex{"hash": options.Hash}, &mnemonic, tx)
	if err != nil {
		return nil, fmt.Errorf("can not find mnemonic by options %v: %w", options, err)
	}
	if mnemonic.Id == 0 {
		return nil, nil
	}

	return &mnemonic, nil
}

func (r Repository) FindOne(id int64, tx driver.Tx) (*Mnemonic, error) {
	mnemonic := Mnemonic{}
	err := r.baseRepository.FindOne(id, &mnemonic, tx)
	if err != nil {
		return nil, fmt.Errorf("can not find mnemonic by id %v: %w", id, err)
	}
	if mnemonic.Id == 0 {
		return nil, nil
	}

	return &mnemonic, nil
}

type FindManyOptions struct {
	UserId int64
}

func (r Repository) FindMany(options FindManyOptions, tx driver.Tx) ([]*Mnemonic, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{
			"id": goqu.Select("mnemonic_id").From("users_mnemonics").Where(goqu.Ex{"user_id": options.UserId}),
		},
		Limit:   50000,
		Offset:  0,
		OrderBy: goqu.I("id").Desc(),
	}, &Mnemonic{}, tx)
	if err != nil {
		return nil, fmt.Errorf("error find user mnemonic: %w", err)
	}

	return items, nil
}

func (r Repository) Save(mnemonic *Mnemonic, tx driver.Tx) error {
	err := r.baseRepository.CreateOrUpdate(mnemonic, tx)
	if err != nil {
		return fmt.Errorf("can not save mnemonic: %s", err)
	}
	err = r.baseRepository.Refresh(mnemonic, tx)
	if err != nil {
		return fmt.Errorf("can not save mnemonic: %s", err)
	}

	return nil
}
