package profile

import "nexus-wallet/internal/app_enum"

type GetInput struct {
	UserId int64
}

type GetOutputTelegram struct {
	FirstName   string
	LastName    string
	Username    string
	ImageSource string
}

type GetOutput struct {
	TelegramProfile GetOutputTelegram
	Type            app_enum.ProfileType
	Language        app_enum.Language
}
