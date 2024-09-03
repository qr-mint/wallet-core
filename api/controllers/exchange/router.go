package exchange

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/middlewares"
)

type Router struct {
	exchangeController    *ExchangeController
	accessTokenValidator  *middlewares.AccessTokenValidator
	mnemonicHashValidator *middlewares.MnemonicHashValidator
	rateLimiter           *middlewares.RateLimiter
}

func NewRouter(
	exchangeController *ExchangeController,
	accessTokenValidator *middlewares.AccessTokenValidator,
	mnemonicHashValidator *middlewares.MnemonicHashValidator,
	rateLimiter *middlewares.RateLimiter,
) Router {
	return Router{
		exchangeController:    exchangeController,
		accessTokenValidator:  accessTokenValidator,
		mnemonicHashValidator: mnemonicHashValidator,
		rateLimiter:           rateLimiter,
	}
}

func (r Router) SetRoutes(group *gin.RouterGroup) {
	routes := group.Group("/exchange")
	{
		routes.Use(r.rateLimiter.CheckLimit())
		routes.Use(r.accessTokenValidator.ValidateAccessToken())
		routes.Use(r.mnemonicHashValidator.ValidateMnemonicHash())
		routes.GET("/list", r.exchangeController.List)
		routes.GET("/address", r.exchangeController.GetAddress)
		routes.GET("/limits", r.exchangeController.GetLimits)
		routes.GET("/amount", r.exchangeController.GetAmount)
	}
}
