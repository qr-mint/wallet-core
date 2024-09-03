package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/auth"
	"strings"
)

type AccessTokenValidator struct {
	service      *auth.Service
	errorHandler *error_handler.HttpErrorHandler
}

func NewAccessTokenValidator(
	service *auth.Service,
	errorHandler *error_handler.HttpErrorHandler,
) *AccessTokenValidator {
	return &AccessTokenValidator{
		service:      service,
		errorHandler: errorHandler,
	}
}

func (v AccessTokenValidator) ValidateAccessToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString, err := v.getTokenFromContext(context)
		if err != nil {
			v.errorHandler.Handle(context, app_error.UnauthorizedError(err))
			context.Abort()
			return
		}

		tokenData, checkErr := v.service.CheckAccessToken(auth.CheckAccessTokenInput{AccessToken: tokenString})
		if checkErr != nil {
			v.errorHandler.Handle(context, checkErr)
			context.Abort()
			return
		}
		if tokenData.UserId == 0 {
			v.errorHandler.Handle(context, app_error.UnauthorizedError(err))
			context.Abort()
			return
		}

		context.Set("userId", tokenData.UserId)
		context.Next()
	}
}

func (v AccessTokenValidator) getTokenFromContext(c *gin.Context) (string, error) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("token must be not empty")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("token must be bearer token")
	}

	return headerParts[1], nil
}
