package notification

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/pkg/repository"
	"time"
)

type GlobalNotification struct {
	Id        int64     `primary:"true" must_generate:"true" db:"id"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (GlobalNotification) GetTableName() string {
	return "global_notifications"
}

func (GlobalNotification) Clear() repository.Model {
	return &GlobalNotification{}
}

type GlobalNotificationTranslation struct {
	Id           int64             `primary:"true" must_generate:"true" db:"id"`
	Text         string            `db:"text"`
	ImagePath    *string           `db:"image_path"`
	LanguageCode app_enum.Language `db:"language"`
}

func (GlobalNotificationTranslation) GetTableName() string {
	return "global_notification_translations"
}

func (GlobalNotificationTranslation) Clear() repository.Model {
	return &GlobalNotificationTranslation{}
}

type GlobalProcessedNotification struct {
	Id             int64 `primary:"true" must_generate:"true" db:"id"`
	NotificationId int64 `db:"notification_id"`
	UserId         int64 `db:"user_id"`
}

func (GlobalProcessedNotification) GetTableName() string {
	return "global_processed_notifications"
}

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
	Id           int64             `primary:"true" must_generate:"true" db:"id"`
	Text         string            `db:"text"`
	ImagePath    *string           `db:"image_path"`
	LanguageCode app_enum.Language `db:"language"`
}

func (PersonalNotificationTranslation) GetTableName() string {
	return "personal_notification_translations"
}

func (PersonalNotificationTranslation) Clear() repository.Model {
	return &PersonalNotificationTranslation{}
}
