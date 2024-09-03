package estimate

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/modules/coin"
)

type EstimateController struct {
	coinService     *coin.Service
	responseFactory *response.ResponseFactory
	errorHandler    *error_handler.HttpErrorHandler
}

func NewEstimateController(
	coinService *coin.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *EstimateController {
	return &EstimateController{
		coinService:     coinService,
		responseFactory: responseFactory,
		errorHandler:    errorHandler,
	}
}

func (c EstimateController) GetCryptoPriceInFiat(context *gin.Context) {
	input, err := GetCryptoPriceInFiatRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	data, err := c.coinService.GetCryptoPriceInFiat(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, GetCryptoPriceInFiatResponse{
		PayableAmount: data.PayableAmount,
		Price:         data.Price,
	})
}

func (c EstimateController) GetFiatPriceInCrypto(context *gin.Context) {
	input, err := GetFiatPriceInCryptoRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	data, err := c.coinService.GetFiatPriceInCrypto(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, GetFiatPriceInCryptoResponse{
		PayableAmount: data.PayableAmount,
		Price:         data.Price,
	})
}
