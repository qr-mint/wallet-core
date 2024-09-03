package telegram

import (
	"errors"
	"gitlab.com/golib4/telegram-query-serice/telegram_query"
	"nexus-wallet/internal/app_error"
	user_factory "nexus-wallet/internal/modules/auth/factory/user"
	"nexus-wallet/internal/modules/auth/model/user"
)

type Service struct {
	telegramQueryService *telegram_query.QueryService
	telegramUserFactory  *user_factory.TelegramUserFactory
	userRepository       *user.TelegramUserRepository
}

func NewService(
	telegramQueryService *telegram_query.QueryService,
	telegramUserFactory *user_factory.TelegramUserFactory,
	userRepository *user.TelegramUserRepository,
) *Service {
	return &Service{
		telegramQueryService: telegramQueryService,
		telegramUserFactory:  telegramUserFactory,
		userRepository:       userRepository,
	}
}

func (s Service) Auth(data AuthInput) (*AuthOutput, *app_error.AppError) {
	err := s.telegramQueryService.Validate(data.TelegramQuery)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid telegram query provided"))
	}

	userData, err := s.telegramQueryService.Parse(data.TelegramQuery)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid telegram query provided, can not parse telegram query"))
	}

	userModel := user_factory.CreateData{
		TelegramId:   userData.ID,
		FirstName:    userData.FirstName,
		LastName:     userData.LastName,
		Username:     userData.Username,
		LanguageCode: userData.LanguageCode,
	}

	userId, appError := s.telegramUserFactory.Create(userModel)
	if appError != nil {
		return nil, appError
	}

	return &AuthOutput{UserId: userId}, nil
}
