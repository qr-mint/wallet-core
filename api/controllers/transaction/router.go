package transaction

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/middlewares"
)

type Router struct {
	transactionController *TransactionController
	accessTokenValidator  *middlewares.AccessTokenValidator
	mnemonicHashValidator *middlewares.MnemonicHashValidator
	rateLimiter           *middlewares.RateLimiter
}

func NewRouter(
	transactionController *TransactionController,
	accessTokenValidator *middlewares.AccessTokenValidator,
	mnemonicHashValidator *middlewares.MnemonicHashValidator,
	rateLimiter *middlewares.RateLimiter,
) Router {
	return Router{
		transactionController: transactionController,
		accessTokenValidator:  accessTokenValidator,
		mnemonicHashValidator: mnemonicHashValidator,
		rateLimiter:           rateLimiter,
	}
}

func (r Router) SetRoutes(group *gin.RouterGroup) {
	routes := group.Group("/transactions")
	{
		routes.Use(r.rateLimiter.CheckLimit())
		routes.Use(r.accessTokenValidator.ValidateAccessToken())
		routes.Use(r.mnemonicHashValidator.ValidateMnemonicHash())
		routes.GET("", r.transactionController.List)
	}
}
