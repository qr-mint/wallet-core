package notification

import (
	"database/sql/driver"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"nexus-wallet/pkg/repository"
)

type PersonalNotificationRepository struct {
	baseRepository *repository.BaseRepository
}

func NewPersonalNotificationRepository(baseRepository *repository.BaseRepository) *PersonalNotificationRepository {
	return &PersonalNotificationRepository{baseRepository: baseRepository}
}

type FindAllOptions struct {
	Limit  uint
	Offset uint
}

func (r PersonalNotificationRepository) FindAll(options FindAllOptions, tx driver.Tx) ([]*PersonalNotification, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: nil,
		Limit:      options.Limit,
		Offset:     options.Offset,
		OrderBy:    goqu.I("id").Asc(),
	}, &PersonalNotification{}, tx)
	if err != nil {
		return nil, errors.Errorf("can not get personal notifications in repository: %s", err.Error())
	}
	return items, nil
}

func (r PersonalNotificationRepository) Delete(notification *PersonalNotification, tx driver.Tx) error {
	err := r.baseRepository.DeleteBy(goqu.Ex{"id": notification.Id}, notification, tx)
	if err != nil {
		return errors.Errorf("can not delete personal notification in repository: %s", err.Error())
	}

	return nil
}
