package deposit

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/middlewares"
)

type Router struct {
	depositController     *DepositController
	accessTokenValidator  *middlewares.AccessTokenValidator
	mnemonicHashValidator *middlewares.MnemonicHashValidator
	rateLimiter           *middlewares.RateLimiter
}

func NewRouter(
	depositController *DepositController,
	accessTokenValidator *middlewares.AccessTokenValidator,
	mnemonicHashValidator *middlewares.MnemonicHashValidator,
	rateLimiter *middlewares.RateLimiter,
) Router {
	return Router{
		depositController:     depositController,
		accessTokenValidator:  accessTokenValidator,
		mnemonicHashValidator: mnemonicHashValidator,
		rateLimiter:           rateLimiter,
	}
}

func (r Router) SetRoutes(group *gin.RouterGroup) {
	routes := group.Group("/deposit")
	{
		routes.Use(r.rateLimiter.CheckLimit())
		routes.Use(r.accessTokenValidator.ValidateAccessToken())
		routes.Use(r.mnemonicHashValidator.ValidateMnemonicHash())
		routes.GET("/:providerName/redirect", r.depositController.GetRedirectUrl)
		routes.GET("/:providerName/limits", r.depositController.GetLimits)
	}
}
