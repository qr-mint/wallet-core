package transfer

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/transfer"
)

type TransferController struct {
	service         *transfer.Service
	responseFactory *response.ResponseFactory
	errorHandler    *error_handler.HttpErrorHandler
}

func NewTransferController(
	service *transfer.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *TransferController {
	return &TransferController{
		service:         service,
		responseFactory: responseFactory,
		errorHandler:    errorHandler,
	}
}

func (c *TransferController) GetMessage(context *gin.Context) {
	input, appErr := GetMessageRequest{}.createInputFromRequest(context)
	if appErr != nil {
		c.errorHandler.Handle(context, appErr)
		return
	}

	result, appErr := c.service.BuildTransfer(*input)
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

func (c *TransferController) Transfer(context *gin.Context) {
	input, err := TransferRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	transferData, transferError := c.service.Transfer(*input)
	if transferError != nil {
		c.errorHandler.Handle(context, transferError)
		return
	}

	c.responseFactory.Ok(context, TransferResponse{Hash: transferData.Hash})
}
