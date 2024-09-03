package telegram

import (
	"fmt"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/notification/model/user"
	"nexus-wallet/internal/modules/notification/processor"
	"strings"
)

type telegramNotificationProcessor struct {
	telegramUserRepository *user.TelegramUserRepository
	botApi                 *tgbotAPI.BotAPI
}

func NewTelegramProcessor(
	botApi *tgbotAPI.BotAPI,
	telegramUserRepository *user.TelegramUserRepository,
) processor.NotificationProcessor {
	return &telegramNotificationProcessor{
		botApi:                 botApi,
		telegramUserRepository: telegramUserRepository,
	}
}

func (n telegramNotificationProcessor) Notify(input processor.NotifyInput) *app_error.AppError {
	telegramUser, err := n.telegramUserRepository.Find(user.FindOptions{UserId: input.UserId}, nil)
	if err != nil {
		return app_error.InternalError(errors.Errorf("can not do find telegram user: %s", err))
	}
	if telegramUser == nil {
		return app_error.InternalError(errors.Errorf("telegram user with user_id %d not found", input.UserId))
	}

	messageConfig := n.createMessage(input, telegramUser.TelegramId)
	_, err = n.botApi.Send(messageConfig)
	if err != nil {
		if strings.Contains(err.Error(), "Forbidden") {
			return nil
		}

		return app_error.InternalError(fmt.Errorf("can not send message to bot: %s", err))
	}

	return nil
}

func (n telegramNotificationProcessor) createMessage(input processor.NotifyInput, telegramId int64) tgbotAPI.Chattable {
	if input.ImagePath == nil {
		return n.createTextMessage(input.Text, telegramId)
	}

	return n.createTextMessageWithImage(*input.ImagePath, input.Text, telegramId)
}

func (n telegramNotificationProcessor) createTextMessage(message string, telegramId int64) tgbotAPI.Chattable {
	messageConfig := tgbotAPI.NewMessage(telegramId, message)
	messageConfig.ParseMode = tgbotAPI.ModeHTML
	return messageConfig
}

func (n telegramNotificationProcessor) createTextMessageWithImage(imagePath string, message string, telegramId int64) tgbotAPI.Chattable {
	messageConfig := tgbotAPI.NewPhoto(telegramId, tgbotAPI.FilePath(imagePath))
	messageConfig.Caption = message
	messageConfig.ParseMode = tgbotAPI.ModeHTML
	return messageConfig
}
