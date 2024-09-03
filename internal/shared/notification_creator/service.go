package notification_creator

import (
	"database/sql/driver"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/shared/notification_creator/model"
	"nexus-wallet/pkg/transaction"
)

type Creator struct {
	repository         *model.Repository
	transactionManager transaction.Manager
}

func NewCreator(
	repository *model.Repository,
	transactionManager transaction.Manager,
) *Creator {
	return &Creator{
		repository:         repository,
		transactionManager: transactionManager,
	}
}

func (s Creator) Create(input CreateInput) *app_error.AppError {
	var translations []*model.PersonalNotificationTranslation
	for language, text := range input.Texts {
		translations = append(translations, &model.PersonalNotificationTranslation{
			Text:         text,
			ImagePath:    input.ImagePath,
			LanguageCode: language,
		})
	}
	err := s.transactionManager.WithTransaction(func(tx driver.Tx) error {
		err := s.repository.Save(&model.PersonalNotification{UserId: input.UserId}, translations, tx)
		if err != nil {
			return errors.Errorf("failed to save notification: %s", err)
		}

		return nil
	})
	if err != nil {
		return app_error.InternalError(errors.Errorf("failed to save notifications in transaction: %s", err))
	}

	return nil
}
