package model

import (
	"database/sql/driver"
	"fmt"
	"nexus-wallet/pkg/repository"
)

type Repository struct {
	baseRepository *repository.BaseRepository
}

func NewRepository(baseRepository *repository.BaseRepository) *Repository {
	return &Repository{baseRepository: baseRepository}
}

func (r Repository) Save(notification *PersonalNotification, translations []*PersonalNotificationTranslation, tx driver.Tx) error {
	err := r.baseRepository.CreateOrUpdate(notification, tx)
	if err != nil {
		return fmt.Errorf("can not save personal notification: %s", err)
	}
	for _, translation := range translations {
		translation.NotificationId = notification.Id
		err := r.baseRepository.CreateOrUpdate(translation, tx)
		if err != nil {
			return fmt.Errorf("can not save personal notification translation: %s", err)
		}
	}

	return nil

}
