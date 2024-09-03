package notification

import (
	"database/sql/driver"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type PersonalNotificationTranslationRepository struct {
	baseRepository *repository.BaseRepository
}

func NewPersonalNotificationTranslationRepository(baseRepository *repository.BaseRepository) *PersonalNotificationTranslationRepository {
	return &PersonalNotificationTranslationRepository{baseRepository: baseRepository}
}

type FindPersonalTranslationsOptions struct {
	NotificationId int64
}

func (r PersonalNotificationTranslationRepository) FindTranslations(options FindPersonalTranslationsOptions, tx driver.Tx) (map[app_enum.Language]PersonalNotificationTranslation, error) {
	translations, err := repository.FindManyBy(r.baseRepository, repository.FindManyByOptions{
		Expression: goqu.Ex{"notification_id": options.NotificationId},
		Limit:      50000,
		Offset:     0,
		OrderBy:    nil,
	}, &PersonalNotificationTranslation{}, tx)
	if err != nil {
		return nil, errors.Errorf("can not get personal notifications translations in repository: %s", err.Error())
	}

	translationMap := make(map[app_enum.Language]PersonalNotificationTranslation)
	for _, translation := range translations {
		translationMap[translation.LanguageCode] = *translation
	}
	return translationMap, nil
}
