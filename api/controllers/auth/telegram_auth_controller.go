package auth

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/modules/auth"
)

type TelegramAuthController struct {
	authService     *auth.Service
	errorHandler    *error_handler.HttpErrorHandler
	responseFactory *response.ResponseFactory
}

func NewTelegramAuthController(
	authService *auth.Service,
	errorHandler *error_handler.HttpErrorHandler,
	responseFactory *response.ResponseFactory,
) *TelegramAuthController {
	return &TelegramAuthController{
		authService:     authService,
		errorHandler:    errorHandler,
		responseFactory: responseFactory,
	}
}

func (c *TelegramAuthController) Auth(context *gin.Context) {
	input, err := AuthRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}
	data, err := c.authService.AuthenticateThroughTelegram(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}
	c.responseFactory.Ok(context, AuthResponse{AccessToken: data.AccessToken, RefreshToken: data.RefreshToken})
}
