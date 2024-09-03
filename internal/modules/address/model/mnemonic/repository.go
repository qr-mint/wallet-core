package mnemonic

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

func (r Repository) FindOne(Id int64, tx driver.Tx) (*Mnemonic, error) {
	mnemonic := Mnemonic{}
	err := r.baseRepository.FindOne(Id, &mnemonic, tx)
	if err != nil {
		return nil, fmt.Errorf("can not find mnemonic by id %v: %w", Id, err)
	}
	if mnemonic.Id == 0 {
		return nil, nil
	}

	return &mnemonic, nil
}

type IsAssignedToUserOptions struct {
	UserId     int64
	MnemonicId int64
}

func (r Repository) IsAssignedToUser(options IsAssignedToUserOptions, tx driver.Tx) (bool, error) {
	usersMnemonics := UsersMnemonics{}
	err := r.baseRepository.FindOneBy(goqu.Ex{"user_id": options.UserId, "mnemonic_id": options.MnemonicId}, &usersMnemonics, tx)
	if err != nil {
		return false, fmt.Errorf("can not check exists mnemonic by options %v: %w", options, err)
	}

	return usersMnemonics.Id != 0, nil
}

func (r Repository) Save(mnemonic *Mnemonic, userId int64, tx driver.Tx) error {
	err := r.baseRepository.CreateOrUpdate(mnemonic, tx)
	if err != nil {
		return fmt.Errorf("can not save mnemonic: %s", err)
	}
	err = r.baseRepository.Refresh(mnemonic, tx)
	if err != nil {
		return fmt.Errorf("can not refresh mnemonic: %s", err)
	}
	err = r.AssignUser(mnemonic, userId, tx)
	if err != nil {
		return fmt.Errorf("can not assign user to mnemonic: %s", err)
	}
	return nil
}

func (r Repository) AssignUser(mnemonic *Mnemonic, userId int64, tx driver.Tx) error {
	usersMnemonics := UsersMnemonics{UserId: userId, MnemonicId: mnemonic.Id}
	err := r.baseRepository.CreateOrUpdate(&usersMnemonics, tx)
	if err != nil {
		return fmt.Errorf("can not create users mnemonics: %s", err)
	}
	err = r.baseRepository.Refresh(&usersMnemonics, tx)
	if err != nil {
		return fmt.Errorf("can not refresh users mnemonics: %s", err)
	}

	return nil
}
