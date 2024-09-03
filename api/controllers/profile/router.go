package profile

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/middlewares"
)

type Router struct {
	profileController    *ProfileController
	accessTokenValidator *middlewares.AccessTokenValidator
	rateLimiter          *middlewares.RateLimiter
}

func NewRouter(
	profileController *ProfileController,
	accessTokenValidator *middlewares.AccessTokenValidator,
	rateLimiter *middlewares.RateLimiter,
) Router {
	return Router{
		profileController:    profileController,
		accessTokenValidator: accessTokenValidator,
		rateLimiter:          rateLimiter,
	}
}

func (r Router) SetRoutes(group *gin.RouterGroup) {
	routes := group.Group("/profile")
	{
		routes.Use(r.rateLimiter.CheckLimit())
		routes.Use(r.accessTokenValidator.ValidateAccessToken())
		routes.GET("", r.profileController.GetProfile)
	}
}
