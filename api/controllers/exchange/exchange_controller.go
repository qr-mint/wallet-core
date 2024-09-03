package exchange

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/modules/exchange"
)

type ExchangeController struct {
	service         *exchange.Service
	responseFactory *response.ResponseFactory
	errorHandler    *error_handler.HttpErrorHandler
}

func NewExchangeController(
	service *exchange.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *ExchangeController {
	return &ExchangeController{
		service:         service,
		responseFactory: responseFactory,
		errorHandler:    errorHandler,
	}
}

func (c ExchangeController) List(context *gin.Context) {
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

func (c ExchangeController) GetAddress(context *gin.Context) {
	input, err := GetAddressForTransferRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	output, err := c.service.ProvideAddressForTransfer(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, GetAddressForTransferResponse{}.fillFromOutput(*output))
}

func (c ExchangeController) GetAmount(context *gin.Context) {
	input, err := GetExchangeAmountRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	output, err := c.service.GetExchangeAmount(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, GetExchangeAmountResponse{}.fillFromOutput(*output))
}

func (c ExchangeController) GetLimits(context *gin.Context) {
	input, err := GetLimitsRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	output, err := c.service.GetLimits(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, GetLimitsResponse{}.fillFromOutput(*output))
}
