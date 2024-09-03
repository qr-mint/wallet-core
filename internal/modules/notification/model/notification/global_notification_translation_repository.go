package notification

import (
	"database/sql/driver"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type GlobalNotificationTranslationRepository struct {
	baseRepository *repository.BaseRepository
}

func NewGlobalNotificationTranslationRepository(baseRepository *repository.BaseRepository) *GlobalNotificationTranslationRepository {
	return &GlobalNotificationTranslationRepository{baseRepository: baseRepository}
}

type FindTranslationsOptions struct {
	NotificationId int64
}

func (r GlobalNotificationTranslationRepository) FindTranslations(options FindTranslationsOptions, tx driver.Tx) (map[app_enum.Language]GlobalNotificationTranslation, error) {
	translations, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{"notification_id": options.NotificationId},
		Limit:      50000,
		Offset:     0,
		OrderBy:    nil,
	}, &GlobalNotificationTranslation{}, tx)
	if err != nil {
		return nil, errors.Errorf("can not get global notifications translations in repository: %s", err.Error())
	}

	translationMap := make(map[app_enum.Language]GlobalNotificationTranslation)
	for _, translation := range translations {
		translationMap[translation.LanguageCode] = *translation
	}
	return translationMap, nil
}
