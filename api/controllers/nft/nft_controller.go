package nft

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/nft"
)

type NftController struct {
	responseFactory *response.ResponseFactory
	errorHandler    *error_handler.HttpErrorHandler
	service         *nft.Service
}

func NewNftController(
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
	service *nft.Service,
) *NftController {
	return &NftController{
		responseFactory: responseFactory,
		errorHandler:    errorHandler,
		service:         service,
	}
}

func (c NftController) List(context *gin.Context) {
	input, err := ListRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	output, err := c.service.List(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, ListResponse{}.fillFromOutput(*output))
}

func (c NftController) Get(context *gin.Context) {
	input, err := GetRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	output, err := c.service.Get(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, GetResponse{}.fillFromModel(*output))
}

func (c NftController) GetMessage(context *gin.Context) {
	input, appErr := GetMessageRequest{}.createInputFromRequest(context)
	if appErr != nil {
		c.errorHandler.Handle(context, appErr)
		return
	}

	result, appErr := c.service.BuildSendMessage(*input)
	if appErr != nil {
		c.errorHandler.Handle(context, appErr)
		return
	}

	messageResponse, err := GetMessageResponse{}.fillFromMessage(result.Message)
	if err != nil {
		c.errorHandler.Handle(context, app_error.InternalError(err))
		return
	}

	c.responseFactory.Ok(context, messageResponse)

}

func (c NftController) Send(context *gin.Context) {
	input, err := SendRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	sendData, sendErr := c.service.Send(*input)
	if sendErr != nil {
		c.errorHandler.Handle(context, sendErr)
		return
	}

	c.responseFactory.Ok(context, SendResponse{Hash: sendData.Hash})

}
