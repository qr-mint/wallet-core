package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/mnemonic"
)

type MnemonicHashValidator struct {
	service      *mnemonic.Service
	errorHandler *error_handler.HttpErrorHandler
}

func NewMnemonicHashValidator(
	service *mnemonic.Service,
	errorHandler *error_handler.HttpErrorHandler,
) *MnemonicHashValidator {
	return &MnemonicHashValidator{service: service, errorHandler: errorHandler}
}

func (v MnemonicHashValidator) ValidateMnemonicHash() gin.HandlerFunc {
	return func(context *gin.Context) {
		hash := v.getHeaderFromContext(context)
		if hash == "" {
			v.errorHandler.Handle(context, app_error.InvalidDataError(errors.New("Mnemonic-Hash header must be not empty")))
			context.Abort()
			return
		}

		mnemonicData, getErr := v.service.GetMnemonicId(mnemonic.GetMnemonicIdInput{Hash: hash})
		if getErr != nil {
			v.errorHandler.Handle(context, getErr)
			context.Abort()
			return
		}
		if mnemonicData.Id == 0 {
			v.errorHandler.Handle(context, app_error.InvalidDataError(errors.New("invalid Mnemonic-Hash provided")))
			context.Abort()
			return
		}

		context.Set("mnemonicId", mnemonicData.Id)
		context.Next()
	}
}

func (v MnemonicHashValidator) getHeaderFromContext(c *gin.Context) string {
	return c.GetHeader("Mnemonic-Hash")
}
