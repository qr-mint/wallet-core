package profile

import (
	"errors"
	"fmt"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/profile/model"
)

type Service struct {
	repository *model.Repository
}

func NewService(repository *model.Repository) *Service {
	return &Service{repository: repository}
}

func (s Service) Get(input GetInput) (*GetOutput, *app_error.AppError) {
	profile, err := s.repository.Find(model.FindOptions{UserId: input.UserId}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("find profile error: %w", err))
	}
	if profile == nil {
		return nil, app_error.InternalError(fmt.Errorf("find profile by user id %d resutned nil result", input.UserId))
	}
	if profile.GetProfile().Type == app_enum.TelegramProfileType {
		telegramProfile, ok := profile.(*model.TelegramProfile)
		if !ok {
			return nil, app_error.InternalError(errors.New("telegram profile type error, invalid casting"))
		}

		return &GetOutput{
			TelegramProfile: GetOutputTelegram{
				FirstName:   telegramProfile.FirstName,
				LastName:    telegramProfile.LastName,
				Username:    telegramProfile.Username,
				ImageSource: telegramProfile.ImageSource,
			},
			Type:     telegramProfile.Profile.Type,
			Language: telegramProfile.Profile.Language,
		}, nil
	}

	return nil, app_error.InternalError(fmt.Errorf("unsupportable profile type: %s", profile.GetProfile().Type))
}
