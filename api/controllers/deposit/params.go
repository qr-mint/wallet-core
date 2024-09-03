package deposit

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/deposit"
	"nexus-wallet/internal/modules/deposit/enum"
	"strconv"
)

type GetRedirectUrlRequest struct {
}

func (GetRedirectUrlRequest) createInputFromRequest(context *gin.Context) (*deposit.ProvideRedirectUrlInput, *app_error.AppError) {
	providerName := enum.ToProviderName(context.Param("providerName"))
	if providerName == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid providerName given"))
	}
	currency := app_enum.ToCurrency(context.GetHeader("Currency-Code"))
	if currency == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid currency given"))
	}
	amount, err := strconv.ParseFloat(context.Request.URL.Query().Get("amount"), 64)
	if err != nil || amount <= 0 {
		return nil, app_error.InvalidDataError(errors.New("invalid amount given"))
	}
	addressCoinId, err := strconv.Atoi(context.Request.URL.Query().Get("address_coin_id"))
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `address_coin_id` given"))
	}
	mnemonicId := context.GetInt64("mnemonicId")

	return &deposit.ProvideRedirectUrlInput{
		MnemonicId:    mnemonicId,
		Amount:        amount,
		AddressCoinId: int64(addressCoinId),
		FiatCurrency:  *currency,
		ProviderName:  *providerName,
	}, nil
}

type GetRedirectUrlResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

type GetLimitsRequest struct {
}

func (GetLimitsRequest) createInputFromRequest(context *gin.Context) (*deposit.GetLimitsInput, *app_error.AppError) {
	providerName := enum.ToProviderName(context.Param("providerName"))
	if providerName == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid providerName given"))
	}
	currency := app_enum.ToCurrency(context.GetHeader("Currency-Code"))
	if currency == nil {
		return nil, app_error.InvalidDataError(errors.New("invalid currency given"))
	}
	addressCoinId, err := strconv.Atoi(context.Request.URL.Query().Get("address_coin_id"))
	if err != nil {
		return nil, app_error.InvalidDataError(errors.New("invalid `address_coin_id` given"))
	}
	mnemonicId := context.GetInt64("mnemonicId")

	return &deposit.GetLimitsInput{
		MnemonicId:    mnemonicId,
		AddressCoinId: int64(addressCoinId),
		FiatCurrency:  *currency,
		ProviderName:  *providerName,
	}, nil
}

type GetLimitsResponse struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

func (GetLimitsResponse) fillFromOutput(output deposit.GetLimitsOutput) GetLimitsResponse {
	return GetLimitsResponse{
		Min: fmt.Sprintf("%.2f", output.Min),
		Max: fmt.Sprintf("%.2f", output.Max),
	}
}
