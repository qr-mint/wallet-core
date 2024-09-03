package deposit

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/modules/deposit"
)

type DepositController struct {
	service         *deposit.Service
	responseFactory *response.ResponseFactory
	errorHandler    *error_handler.HttpErrorHandler
}

func NewDepositController(
	service *deposit.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *DepositController {
	return &DepositController{
		service:         service,
		responseFactory: responseFactory,
		errorHandler:    errorHandler,
	}
}

func (c DepositController) GetRedirectUrl(context *gin.Context) {
	input, err := GetRedirectUrlRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}
	redirectUrl, appErr := c.service.ProvideRedirectUrl(*input)
	if appErr != nil {
		c.errorHandler.Handle(context, appErr)
		return
	}

	c.responseFactory.Ok(context, GetRedirectUrlResponse{RedirectUrl: redirectUrl.Url})
}

func (c DepositController) GetLimits(context *gin.Context) {
	input, err := GetLimitsRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}
	limits, appErr := c.service.GetLimits(*input)
	if appErr != nil {
		c.errorHandler.Handle(context, appErr)
		return
	}

	c.responseFactory.Ok(context, GetLimitsResponse{}.fillFromOutput(*limits))
}
