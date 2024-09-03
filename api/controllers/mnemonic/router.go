package mnemonic

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/middlewares"
)

type Router struct {
	mnemonicController    *MnemonicController
	accessTokenValidator  *middlewares.AccessTokenValidator
	mnemonicHashValidator *middlewares.MnemonicHashValidator
	rateLimiter           *middlewares.RateLimiter
}

func NewRouter(
	mnemonicController *MnemonicController,
	accessTokenValidator *middlewares.AccessTokenValidator,
	mnemonicHashValidator *middlewares.MnemonicHashValidator,
	rateLimiter *middlewares.RateLimiter,
) Router {
	return Router{
		mnemonicController:    mnemonicController,
		accessTokenValidator:  accessTokenValidator,
		mnemonicHashValidator: mnemonicHashValidator,
		rateLimiter:           rateLimiter,
	}
}

func (r Router) SetRoutes(group *gin.RouterGroup) {
	routes := group.Group("/mnemonic")
	{
		routes.Use(r.rateLimiter.CheckLimit())
		routes.Use(r.accessTokenValidator.ValidateAccessToken())
		routes.POST("/validate", r.mnemonicController.ValidateHash)
		routes.POST("/generate", r.mnemonicController.Generate)

		routes.Use(r.accessTokenValidator.ValidateAccessToken())
		routes.GET("/names", r.mnemonicController.GetNames)
		routes.Use(r.mnemonicHashValidator.ValidateMnemonicHash())
		routes.PUT("/name", r.mnemonicController.UpdateName)
	}
}
