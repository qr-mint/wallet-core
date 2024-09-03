package notification

import (
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/notification/model/notification"
	"nexus-wallet/internal/modules/notification/model/profile"
	"nexus-wallet/internal/modules/notification/processor"
)

type Params struct {
	IsFake bool
}

type Service struct {
	processors                                map[app_enum.ProfileType]processor.NotificationProcessor
	profileRepository                         *profile.Repository
	globalNotificationRepository              *notification.GlobalNotificationRepository
	globalNotificationTranslationRepository   *notification.GlobalNotificationTranslationRepository
	globalProcessedNotificationRepository     *notification.GlobalProcessedNotificationRepository
	personalNotificationRepository            *notification.PersonalNotificationRepository
	personalNotificationTranslationRepository *notification.PersonalNotificationTranslationRepository
	params                                    Params
}

func NewService(
	processors map[app_enum.ProfileType]processor.NotificationProcessor,
	profileRepository *profile.Repository,
	globalNotificationRepository *notification.GlobalNotificationRepository,
	globalNotificationTranslationRepository *notification.GlobalNotificationTranslationRepository,
	globalProcessedNotificationRepository *notification.GlobalProcessedNotificationRepository,
	personalNotificationRepository *notification.PersonalNotificationRepository,
	personalNotificationTranslationRepository *notification.PersonalNotificationTranslationRepository,
	params Params,
) *Service {
	return &Service{
		processors:                                processors,
		profileRepository:                         profileRepository,
		globalNotificationRepository:              globalNotificationRepository,
		globalNotificationTranslationRepository:   globalNotificationTranslationRepository,
		globalProcessedNotificationRepository:     globalProcessedNotificationRepository,
		personalNotificationRepository:            personalNotificationRepository,
		personalNotificationTranslationRepository: personalNotificationTranslationRepository,
		params: params,
	}
}

func (s Service) PersonalNotify() *app_error.AppError {
	var personalNotifications []*notification.PersonalNotification
	var err error
	var offset uint = 0
	var limit uint = 1000
	for {
		personalNotifications, err = s.personalNotificationRepository.FindAll(notification.FindAllOptions{
			Limit:  limit,
			Offset: offset,
		}, nil)
		if err != nil {
			return app_error.InternalError(errors.Errorf("can not get personal notifications: %s", err))
		}
		if len(personalNotifications) == 0 {
			break
		}

		for _, personalNotificationData := range personalNotifications {
			translationsOptions := notification.FindPersonalTranslationsOptions{NotificationId: personalNotificationData.Id}
			translations, err := s.personalNotificationTranslationRepository.FindTranslations(translationsOptions, nil)
			if err != nil {
				return app_error.InternalError(errors.Errorf("can not get personal notification %d translations: %s", personalNotificationData.Id, err))
			}

			profileData, err := s.profileRepository.Find(profile.FindOptions{UserId: personalNotificationData.UserId}, nil)
			if err != nil {
				return app_error.InternalError(errors.Errorf("can not get personal notification %d profile: %s", personalNotificationData.Id, err))
			}
			if profileData == nil {
				return app_error.InternalError(errors.Errorf("can not find personal notification %d profile by user_id: %d", personalNotificationData.Id, personalNotificationData.UserId))
			}

			appErr := s.sendNotificationToUser(*profileData, translations[profileData.Language].ImagePath, translations[profileData.Language].Text)
			if appErr != nil {
				return appErr
			}

			err = s.personalNotificationRepository.Delete(personalNotificationData, nil)
			if err != nil {
				return app_error.InternalError(errors.Errorf("can not delete personal notification %s", err))
			}
		}

		offset = offset + limit
	}

	return nil
}

func (s Service) GlobalNotify() *app_error.AppError {
	globalNotifications, err := s.globalNotificationRepository.FindAll(nil)
	if err != nil {
		return app_error.InternalError(errors.Errorf("can not get globalNotifications: %s", err))
	}
	if len(globalNotifications) == 0 {
		return nil
	}

	for _, globalNotificationData := range globalNotifications {
		translationsOptions := notification.FindTranslationsOptions{NotificationId: globalNotificationData.Id}
		translations, err := s.globalNotificationTranslationRepository.FindTranslations(translationsOptions, nil)
		if err != nil {
			return app_error.InternalError(errors.Errorf("can not get global notification %d translations: %s", globalNotificationData.Id, err))
		}

		appErr := s.startGlobalSend(*globalNotificationData, translations)
		if appErr != nil {
			return appErr
		}

		err = s.globalNotificationRepository.Delete(globalNotificationData, nil)
		if err != nil {
			return app_error.InternalError(errors.Errorf("can not delete global notification %d: %s", globalNotificationData.Id, err))
		}
	}

	return nil
}

func (s Service) startGlobalSend(notificationData notification.GlobalNotification, translations map[app_enum.Language]notification.GlobalNotificationTranslation) *app_error.AppError {
	var limit uint = 1000
	var offset uint = 0
	for {
		profiles, err := s.profileRepository.FindAll(profile.FindAllOptions{Limit: limit, Offset: offset}, nil)
		if err != nil {
			return app_error.InternalError(errors.Errorf("can not get profiles: %s", err))
		}
		if len(profiles) == 0 {
			break
		}

		for _, profileData := range profiles {
			processedOptions := notification.FindOptions{UserId: profileData.UserId, NotificationId: notificationData.Id}
			isAlreadyProcessed, err := s.globalProcessedNotificationRepository.Find(processedOptions, nil)
			if err != nil {
				return app_error.InternalError(errors.Errorf("can not check if global notiication is processed: %s", err.Error()))
			}
			if isAlreadyProcessed {
				continue
			}

			appErr := s.sendNotificationToUser(*profileData, translations[profileData.Language].ImagePath, translations[profileData.Language].Text)
			if appErr != nil {
				return appErr
			}

			err = s.globalProcessedNotificationRepository.Save(&notification.GlobalProcessedNotification{NotificationId: notificationData.Id, UserId: profileData.UserId}, nil)
			if err != nil {
				return app_error.InternalError(errors.Errorf("can not set global notification %d processed: %s", notificationData.Id, err.Error()))
			}

		}

		offset = offset + limit
	}

	if offset == 0 {
		return app_error.InternalError(errors.New("empty profiles list returned in global notification send"))
	}

	return nil
}

func (s Service) sendNotificationToUser(profileData profile.Profile, imagePath *string, text string) *app_error.AppError {
	notificationProcessor, processorExists := s.processors[profileData.Type]
	if !processorExists {
		return app_error.InternalError(errors.Errorf("can not find processor for notification type: %s", profileData.Type))
	}
	if s.params.IsFake {
		return nil
	}

	notifyErr := notificationProcessor.Notify(processor.NotifyInput{ImagePath: imagePath, Text: text, UserId: profileData.UserId})
	if notifyErr != nil {
		return notifyErr
	}

	return nil
}
