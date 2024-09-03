package model

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
)

type PersonalNotification struct {
	Id     int64 `primary:"true" must_generate:"true" db:"id"`
	UserId int64 `db:"user_id"`
}

func (PersonalNotification) GetTableName() string {
	return "personal_notifications"
}

func (PersonalNotification) Clear() repository.Model {
	return &PersonalNotification{}
}

type PersonalNotificationTranslation struct {
	Id             int64             `primary:"true" must_generate:"true" db:"id"`
	Text           string            `db:"text"`
	ImagePath      *string           `db:"image_path"`
	LanguageCode   app_enum.Language `db:"language"`
	NotificationId int64             `db:"notification_id"`
}

func (PersonalNotificationTranslation) GetTableName() string {
	return "personal_notification_translations"
}
