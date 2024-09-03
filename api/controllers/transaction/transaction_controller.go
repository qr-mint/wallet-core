package transaction

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/modules/transaction"
)

type TransactionController struct {
	service         *transaction.Service
	responseFactory *response.ResponseFactory
	errorHandler    *error_handler.HttpErrorHandler
}

func NewTransactionController(
	service *transaction.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *TransactionController {
	return &TransactionController{
		service:         service,
		responseFactory: responseFactory,
		errorHandler:    errorHandler,
	}
}

func (c TransactionController) List(context *gin.Context) {
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
