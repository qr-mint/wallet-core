package mnemonic

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/modules/mnemonic"
)

type MnemonicController struct {
	service         *mnemonic.Service
	responseFactory *response.ResponseFactory
	errorHandler    *error_handler.HttpErrorHandler
}

func NewMnemonicController(
	service *mnemonic.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *MnemonicController {
	return &MnemonicController{
		service:         service,
		responseFactory: responseFactory,
		errorHandler:    errorHandler,
	}
}

func (c MnemonicController) ValidateHash(context *gin.Context) {
	input, err := ValidateHashRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	data, err := c.service.GetMnemonicId(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, ValidateHashResponse{IsValid: data.Id != 0})
}

func (c MnemonicController) Generate(context *gin.Context) {
	mnemonicString, err := c.service.Generate(GenerateRequest{}.createInputFromRequest(context))
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}
	c.responseFactory.Ok(context, GenerateResponse{Value: mnemonicString.Mnemonic})
}

func (c MnemonicController) GetNames(context *gin.Context) {
	names, err := c.service.GetNames(GetNamesRequest{}.createInputFromRequest(context))
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, GetNamesResponse{}.fillFromOutput(*names))
}

func (c MnemonicController) UpdateName(context *gin.Context) {
	input, err := UpdateNameRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}
	err = c.service.UpdateName(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, nil)
}
