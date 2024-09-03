package address

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/middlewares"
)

type Router struct {
	addressController     *AddressController
	accessTokenValidator  *middlewares.AccessTokenValidator
	mnemonicHashValidator *middlewares.MnemonicHashValidator
	rateLimiter           *middlewares.RateLimiter
}

func NewRouter(
	addressController *AddressController,
	accessTokenValidator *middlewares.AccessTokenValidator,
	mnemonicHashValidator *middlewares.MnemonicHashValidator,
	rateLimiter *middlewares.RateLimiter,
) Router {
	return Router{
		addressController:     addressController,
		accessTokenValidator:  accessTokenValidator,
		mnemonicHashValidator: mnemonicHashValidator,
		rateLimiter:           rateLimiter,
	}
}

func (r Router) SetRoutes(group *gin.RouterGroup) {
	routes := group.Group("/addresses")
	{
		routes.Use(r.rateLimiter.CheckLimit())
		routes.Use(r.accessTokenValidator.ValidateAccessToken())
		routes.POST("/import", r.addressController.Import)
		routes.Use(r.mnemonicHashValidator.ValidateMnemonicHash())
		routes.GET("/info", r.addressController.Info)
		routes.GET("/coins", r.addressController.ListCoins)
		routes.GET("/coins/:id", r.addressController.Coin)
		routes.POST("/coins/:id/switch", r.addressController.SwitchVisibility)
	}
}
