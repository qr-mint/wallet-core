package simple_swap

import (
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/deposit/provider"
	"nexus-wallet/pkg/simple_swap"
	"strconv"
)

type simpleSwapProvider struct {
	client *simple_swap.Client
}

func NewSimpleSwapProvider(client *simple_swap.Client) provider.RedirectProvider {
	return simpleSwapProvider{client: client}
}

func (p simpleSwapProvider) GetLimits(
	coinName app_enum.CoinName,
	fiatCurrency app_enum.Currency,
) (*provider.Limits, *app_error.AppError) {
	currencyFrom, err := p.transformCurrency(fiatCurrency)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to transform currency from: %v", err))
	}
	currencyTo, err := p.transformCoinName(coinName)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to transform coin name: %s", err))
	}

	ranges, err := p.client.GetRanges(simple_swap.GetRangesRequest{
		Fixed:        false,
		CurrencyFrom: currencyFrom,
		CurrencyTo:   currencyTo,
	})
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to get simple swap range: %s", err))
	}

	minAmount, err := strconv.ParseFloat(ranges.Min, 64)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to parse min amount: %s", err))
	}
	maxAmount, err := strconv.ParseFloat(ranges.Max, 64)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to parse max amount: %s", err))
	}

	return &provider.Limits{
		Min: minAmount,
		Max: maxAmount,
	}, nil
}

func (p simpleSwapProvider) ProvideRedirectUrl(
	coinName app_enum.CoinName,
	fiatCurrency app_enum.Currency,
	amount float64,
	addressTo string,
) (string, *app_error.AppError) {
	currencyFrom, err := p.transformCurrency(fiatCurrency)
	if err != nil {
		return "", app_error.InternalError(errors.Errorf("failed to transform currency from: %v", err))
	}
	currencyTo, err := p.transformCoinName(coinName)
	if err != nil {
		return "", app_error.InternalError(errors.Errorf("failed to transform coin name: %s", err))
	}

	limits, limitsErr := p.GetLimits(coinName, fiatCurrency)
	if limitsErr != nil {
		return "", limitsErr
	}

	if amount < limits.Min || amount > limits.Max {
		return "", app_error.InvalidDataError(errors.Errorf("amount out of range. min: %f, max: %f", limits.Min, limits.Max))
	}

	response, err := p.client.CreateExchange(simple_swap.CreateExchangeRequest{
		Fixed:        false,
		CurrencyFrom: currencyFrom,
		CurrencyTo:   currencyTo,
		Amount:       amount,
		AddressTo:    addressTo,
	})
	if err != nil {
		return "", app_error.InternalError(errors.Errorf("failed to create simple swap exchange in provider: %s", err))
	}

	return response.RedirectUrl, nil
}

func (p simpleSwapProvider) transformCurrency(fiatCurrency app_enum.Currency) (string, error) {
	switch fiatCurrency {
	case app_enum.UsdCurrency:
		return "usd", nil
	case app_enum.EurCurrency:
		return "usd", nil
	case app_enum.RubCurrency:
		return "usd", nil
	}

	return "", errors.Errorf("unknown fiatCurrency provided in simple swap provider: %s", fiatCurrency)
}

func (p simpleSwapProvider) transformCoinName(coinName app_enum.CoinName) (string, error) {
	switch coinName {
	case app_enum.TonCoinName:
		return "ton", nil
	case app_enum.TronCoinName:
		return "trx", nil
	case app_enum.TetherCoinName:
		return "usdttrc20", nil
	}

	return "", errors.Errorf("unknown coinName provided in simple swap provider: %s", coinName)
}
