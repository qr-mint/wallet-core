package notification

import (
	"database/sql/driver"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"nexus-wallet/pkg/repository"
)

type GlobalProcessedNotificationRepository struct {
	baseRepository *repository.BaseRepository
}

func NewGlobalProcessedNotificationRepository(baseRepository *repository.BaseRepository) *GlobalProcessedNotificationRepository {
	return &GlobalProcessedNotificationRepository{baseRepository: baseRepository}
}

type FindOptions struct {
	UserId         int64
	NotificationId int64
}

func (r GlobalProcessedNotificationRepository) Find(options FindOptions, tx driver.Tx) (bool, error) {
	processedNotification := &GlobalProcessedNotification{}
	err := r.baseRepository.FindOneBy(
		goqu.Ex{"user_id": options.UserId, "notification_id": options.NotificationId},
		processedNotification,
		tx,
	)
	if err != nil {
		return false, errors.Errorf("can not get global processedNotification from repository: %s", err.Error())
	}

	return processedNotification.Id != 0, nil
}

func (r GlobalProcessedNotificationRepository) Save(processedNotification *GlobalProcessedNotification, tx driver.Tx) error {
	err := r.baseRepository.CreateOrUpdate(processedNotification, tx)
	if err != nil {
		return errors.Errorf("can not save global processedNotification: %s", err.Error())
	}

	return nil
}
