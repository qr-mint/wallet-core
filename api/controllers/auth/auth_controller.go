package auth

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/modules/auth"
)

type AuthController struct {
	service         *auth.Service
	errorHandler    *error_handler.HttpErrorHandler
	responseFactory *response.ResponseFactory
}

func NewAuthController(service *auth.Service, errorHandler *error_handler.HttpErrorHandler, responseFactory *response.ResponseFactory) *AuthController {
	return &AuthController{
		service:         service,
		errorHandler:    errorHandler,
		responseFactory: responseFactory,
	}
}

func (c AuthController) Refresh(context *gin.Context) {
	input, err := RefreshRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	data, err := c.service.Refresh(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, AuthResponse{
		data.AccessToken,
		data.RefreshToken,
	})
}
