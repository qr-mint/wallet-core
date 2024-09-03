package notification

import (
	"database/sql/driver"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/pkg/repository"
	"time"
)

type GlobalNotificationRepository struct {
	baseRepository *repository.BaseRepository
}

func NewGlobalNotificationRepository(baseRepository *repository.BaseRepository) *GlobalNotificationRepository {
	return &GlobalNotificationRepository{baseRepository: baseRepository}
}

func (r GlobalNotificationRepository) FindAll(tx driver.Tx) ([]*GlobalNotification, error) {
	items, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.I("expires_at").Gt(app_util.TimeToStartOfDay(time.Now())),
		Limit:      50000,
		Offset:     0,
		OrderBy:    nil,
	}, &GlobalNotification{}, tx)
	if err != nil {
		return nil, errors.Errorf("can not get global notifications in repository: %s", err.Error())
	}
	return items, nil
}

func (r GlobalNotificationRepository) Delete(notification *GlobalNotification, tx driver.Tx) error {
	err := r.baseRepository.DeleteBy(goqu.Ex{"id": notification.Id}, notification, tx)
	if err != nil {
		return errors.Errorf("can not delete global notification in repository: %s", err.Error())
	}

	return nil
}
