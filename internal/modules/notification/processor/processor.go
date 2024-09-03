package processor

import "nexus-wallet/internal/app_error"

type NotificationProcessor interface {
	Notify(input NotifyInput) *app_error.AppError
}

type NotifyInput struct {
	ImagePath *string
	Text      string
	UserId    int64
}
