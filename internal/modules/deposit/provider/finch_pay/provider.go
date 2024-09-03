package finch_pay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/deposit/provider"
	"nexus-wallet/pkg/finch_pay"
	"strconv"
)

type Params struct {
	WidgetHost string
	PartnerId  string
	SecretKey  string
}

type finchPayProvider struct {
	params Params
	client *finch_pay.Client
}

func NewFinchPayProvider(params Params, client *finch_pay.Client) provider.RedirectProvider {
	return &finchPayProvider{
		params: params,
		client: client,
	}
}

func (p finchPayProvider) GetLimits(
	coinName app_enum.CoinName,
	fiatCurrency app_enum.Currency,
) (*provider.Limits, *app_error.AppError) {
	currencyFrom, err := p.transformCurrency(fiatCurrency)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to transform currency from: %v", err))
	}

	limits, err := p.client.GetLimits(finch_pay.GetLimitsRequest{Currency: currencyFrom})
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to get fin tech limits: %v", err))
	}

	minAmount, err := strconv.ParseFloat(limits.MinAmount, 64)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to parse min amount: %v", err))
	}
	maxAmount, err := strconv.ParseFloat(limits.MaxAmount, 64)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to parse max amount: %v", err))
	}

	return &provider.Limits{
		Min: minAmount,
		Max: maxAmount,
	}, nil
}

func (p finchPayProvider) ProvideRedirectUrl(
	coinName app_enum.CoinName,
	fiatCurrency app_enum.Currency,
	amount float64,
	addressTo string,
) (string, *app_error.AppError) {
	currencyFrom, err := p.transformCurrency(fiatCurrency)
	if err != nil {
		return "", app_error.InternalError(errors.Errorf("failed to transform currency from: %v", err))
	}
	currencyTo, network, err := p.transformCoinName(coinName)
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

	sig := hmac.New(sha256.New, []byte(p.params.SecretKey))
	sig.Write([]byte(addressTo))

	url := fmt.Sprintf(
		"%s?a=%f&p=%s&c=%s&wallet_address=%s&sign=%s&partner_id=%s",
		p.params.WidgetHost,
		amount,
		currencyFrom,
		currencyTo,
		addressTo,
		hex.EncodeToString(sig.Sum(nil)),
		p.params.PartnerId,
	)

	if network != "" {
		url = fmt.Sprintf("%s&n=%s", url, network)
	}

	return url, nil
}

func (p finchPayProvider) transformCurrency(fiatCurrency app_enum.Currency) (string, error) {
	switch fiatCurrency {
	case app_enum.UsdCurrency:
		return "USD", nil
	case app_enum.EurCurrency:
		return "USD", nil
	case app_enum.RubCurrency:
		return "USD", nil
	}

	return "", errors.Errorf("unknown fiatCurrency provided in simple swap provider: %s", fiatCurrency)
}

func (p finchPayProvider) transformCoinName(coinName app_enum.CoinName) (string, string, error) {
	switch coinName {
	case app_enum.TonCoinName:
		return "TON", "", nil
	case app_enum.TronCoinName:
		return "TRX", "", nil
	case app_enum.TetherCoinName:
		return "USDT", "TRC20", nil
	}

	return "", "", errors.Errorf("unknown coinName provided in simple swap provider: %s", coinName)
}
