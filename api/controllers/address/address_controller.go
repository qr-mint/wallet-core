package address

import (
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/modules/address"
)

type AddressController struct {
	service         *address.Service
	responseFactory *response.ResponseFactory
	errorHandler    *error_handler.HttpErrorHandler
}

func NewAddressController(
	service *address.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *AddressController {
	return &AddressController{
		service:         service,
		responseFactory: responseFactory,
		errorHandler:    errorHandler,
	}
}

func (c AddressController) Import(context *gin.Context) {
	input, err := ImportRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	err = c.service.Import(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, nil)
}

func (c AddressController) Info(context *gin.Context) {
	input, err := InfoRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	data, err := c.service.GetAggregatedInfo(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, InfoResponse{}.fillFromOutput(*data))
}

func (c AddressController) SwitchVisibility(context *gin.Context) {
	input, err := SwitchVisibilityRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	err = c.service.SwitchCoinVisibility(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}
	c.responseFactory.Ok(context, nil)
}

func (c AddressController) Coin(context *gin.Context) {
	input, err := CoinRequest{}.createInputFromRequest(context)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	coinData, err := c.service.GetCoin(*input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, CoinResponse{}.fillFromOutput(*coinData))
}

func (c AddressController) ListCoins(context *gin.Context) {
	input := ListCoinsRequest{}.createInputFromRequest(context)
	list, err := c.service.GetCoinsList(input)
	if err != nil {
		c.errorHandler.Handle(context, err)
		return
	}

	c.responseFactory.Ok(context, ListCoinsResponse{}.fillFromOutput(*list))
}
