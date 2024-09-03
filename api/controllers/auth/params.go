package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/auth"
	"nexus-wallet/internal/modules/auth/telegram"
)

type RefreshRequest struct {
}

func (RefreshRequest) createInputFromRequest(context *gin.Context) (*auth.RefreshInput, *app_error.AppError) {
	refreshToken := context.GetHeader("Refresh-Token")
	if refreshToken == "" {
		return nil, app_error.InvalidDataError(errors.New("Refresh-Token Header is required."))
	}

	return &auth.RefreshInput{RefreshToken: refreshToken}, nil
}

type AuthRequest struct {
}

func (c AuthRequest) createInputFromRequest(context *gin.Context) (*telegram.AuthInput, *app_error.AppError) {
	var body TelegramAuthRequest
	if err := context.BindJSON(&body); err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid json data"))
	}

	return &telegram.AuthInput{
		TelegramQuery: body.TelegramQuery,
	}, nil
}

type TelegramAuthRequest struct {
	TelegramQuery string `json:"telegram_query"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
