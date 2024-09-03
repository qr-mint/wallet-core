package notification_creator

import "nexus-wallet/internal/app_enum"

type CreateInput struct {
	UserId    int64
	Texts     map[app_enum.Language]string
	ImagePath *string
}
