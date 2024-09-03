package auth

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/middlewares"
)

type Router struct {
	authController         *AuthController
	telegramAuthController *TelegramAuthController
	rateLimiter            *middlewares.RateLimiter
}

func NewRouter(
	authController *AuthController,
	telegramAuthController *TelegramAuthController,
	rateLimiter *middlewares.RateLimiter,
) Router {
	return Router{
		authController:         authController,
		telegramAuthController: telegramAuthController,
		rateLimiter:            rateLimiter,
	}
}

func (r Router) SetRoutes(group *gin.RouterGroup) {
	routes := group.Group("/auth")
	{
		routes.Use(r.rateLimiter.CheckLimit())
		routes.POST("/refresh", r.authController.Refresh)
		routes.POST("/telegram", r.telegramAuthController.Auth)
	}
}
