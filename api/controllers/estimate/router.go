package estimate

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/middlewares"
)

type Router struct {
	estimateController   *EstimateController
	accessTokenValidator *middlewares.AccessTokenValidator
	rateLimiter          *middlewares.RateLimiter
}

func NewRouter(
	estimateController *EstimateController,
	accessTokenValidator *middlewares.AccessTokenValidator,
	rateLimiter *middlewares.RateLimiter,
) Router {
	return Router{
		estimateController:   estimateController,
		accessTokenValidator: accessTokenValidator,
		rateLimiter:          rateLimiter,
	}
}

func (r Router) SetRoutes(group *gin.RouterGroup) {
	routes := group.Group("/estimate")
	{
		routes.Use(r.rateLimiter.CheckLimit())
		routes.Use(r.accessTokenValidator.ValidateAccessToken())
		routes.GET("/crypto-in-fiat", r.estimateController.GetCryptoPriceInFiat)
		routes.GET("/fiat-in-crypto", r.estimateController.GetFiatPriceInCrypto)
	}
}
