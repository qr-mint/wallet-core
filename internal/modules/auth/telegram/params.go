package telegram

type AuthInput struct {
	TelegramQuery string
}

type AuthOutput struct {
	UserId int64
}
