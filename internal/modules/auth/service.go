package auth

import (
	"fmt"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/auth/telegram"
	"nexus-wallet/internal/shared/jwt"
)

type Service struct {
	telegramService  *telegram.Service
	tokenPairService *jwt.TokenPairService
}

func NewService(
	telegramService *telegram.Service,
	tokenPairService *jwt.TokenPairService,
) *Service {
	return &Service{
		telegramService:  telegramService,
		tokenPairService: tokenPairService,
	}
}

func (s Service) AuthenticateThroughTelegram(input telegram.AuthInput) (*AuthenticateOutput, *app_error.AppError) {
	result, err := s.telegramService.Auth(input)
	if err != nil {
		return nil, err
	}
	tokenPair, tokenErr := s.tokenPairService.GenerateTokenPair(result.UserId)
	if tokenErr != nil {
		return nil, app_error.InternalError(fmt.Errorf("error generating token pair: %s", tokenErr))
	}

	return &AuthenticateOutput{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s Service) CheckAccessToken(input CheckAccessTokenInput) (*CheckAccessTokenOutput, *app_error.AppError) {
	data, err := s.tokenPairService.CheckAccessToken(input.AccessToken)
	if err != nil {
		return nil, app_error.UnauthorizedError(fmt.Errorf("error checking token: %s", err))
	}
	return &CheckAccessTokenOutput{
		UserId: data.UserId,
	}, nil
}

func (s Service) Refresh(input RefreshInput) (*AuthenticateOutput, *app_error.AppError) {
	data, err := s.tokenPairService.CheckRefreshToken(input.RefreshToken)
	if err != nil {
		return nil, app_error.InvalidDataError(fmt.Errorf("error checking refresh token: %s", err))
	}

	tokenPair, err := s.tokenPairService.GenerateTokenPair(data.UserId)
	if err != nil {
		return nil, app_error.InvalidDataError(fmt.Errorf("can not create token pair: %s", err))
	}

	return &AuthenticateOutput{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
