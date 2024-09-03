package profile

import (
	"database/sql/driver"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	baseRepository *repository.BaseRepository
}

func NewRepository(baseRepository *repository.BaseRepository) *Repository {
	return &Repository{baseRepository}
}

type FindAllOptions struct {
	Limit  uint
	Offset uint
}

func (r Repository) FindAll(options FindAllOptions, tx driver.Tx) ([]*Profile, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Limit:  options.Limit,
		Offset: options.Offset,
	}, &Profile{}, tx)
	if err != nil {
		return nil, errors.Errorf("can not get profiles in repository: %s", err.Error())
	}

	return items, nil
}

type FindOptions struct {
	UserId int64
}

func (r Repository) Find(options FindOptions, tx driver.Tx) (*Profile, error) {
	profile := Profile{}
	err := r.baseRepository.FindOneBy(goqu.Ex{"user_id": options.UserId}, &profile, tx)
	if err != nil {
		return nil, errors.Errorf("can not find profile in repository: %s", err.Error())
	}

	if profile.Id == 0 {
		return nil, nil
	}

	return &profile, nil
}
