package estimate

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/coin"
	"strconv"
)

type GetCryptoPriceInFiatRequest struct {
}

func (GetCryptoPriceInFiatRequest) createInputFromRequest(context *gin.Context) (*coin.GetCryptoPriceInFiatInput, *app_error.AppError) {
	fiatCurrency := app_enum.ToCurrency(context.GetHeader("Currency-Code"))
	if fiatCurrency == nil {
		return nil, app_error.InvalidDataError(errors.New("Invalid Currency-Code provided"))
	}
	coinName := app_enum.ToCoinName(context.Request.URL.Query().Get("coin_name"))
	if coinName == nil {
		return nil, app_error.InvalidDataError(errors.New("Invalid `coin_name` provided"))
	}
	network := app_enum.ToNetwork(context.Request.URL.Query().Get("network"))
	if network == nil {
		return nil, app_error.InvalidDataError(errors.New("Invalid network provided"))
	}
	amount, err := strconv.ParseFloat(context.Request.URL.Query().Get("crypto_amount"), 64)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `crypto_amount` given. param must be numeric and required"))
	}

	return &coin.GetCryptoPriceInFiatInput{
		Network:      *network,
		CoinName:     *coinName,
		FiatCurrency: *fiatCurrency,
		Amount:       amount,
	}, nil
}

type GetCryptoPriceInFiatResponse struct {
	PayableAmount float64 `json:"payable_amount"`
	Price         float64 `json:"price"`
}

type GetFiatPriceInCryptoRequest struct {
}

func (GetFiatPriceInCryptoRequest) createInputFromRequest(context *gin.Context) (*coin.GetFiatPriceInCryptoInput, *app_error.AppError) {
	fiatCurrency := app_enum.ToCurrency(context.GetHeader("Currency-Code"))
	if fiatCurrency == nil {
		return nil, app_error.InvalidDataError(errors.New("Invalid Currency-Code provided"))
	}
	coinName := app_enum.ToCoinName(context.Request.URL.Query().Get("coin_name"))
	if coinName == nil {
		return nil, app_error.InvalidDataError(errors.New("Invalid `coin_name` provided"))
	}
	network := app_enum.ToNetwork(context.Request.URL.Query().Get("network"))
	if network == nil {
		return nil, app_error.InvalidDataError(errors.New("Invalid network provided"))
	}
	amount, err := strconv.ParseFloat(context.Request.URL.Query().Get("fiat_amount"), 64)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `fiat_amount` given. param must be numeric and required"))
	}

	return &coin.GetFiatPriceInCryptoInput{
		Network:      *network,
		CoinName:     *coinName,
		FiatCurrency: *fiatCurrency,
		Amount:       amount,
	}, nil
}

type GetFiatPriceInCryptoResponse struct {
	PayableAmount float64 `json:"payable_amount"`
	Price         float64 `json:"price"`
}
