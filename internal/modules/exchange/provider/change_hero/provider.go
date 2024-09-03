package change_hero

import (
	"fmt"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/internal/modules/exchange/model/coin"
	"nexus-wallet/internal/modules/exchange/provider"
	"nexus-wallet/pkg/change_hero"
	"strconv"
)

type changeHeroProvider struct {
	client          *change_hero.Client
	priceRepository *coin.PriceRepository
}

func NewChangeHeroProvider(client *change_hero.Client, priceRepository *coin.PriceRepository) provider.TransferableProvider {
	return changeHeroProvider{
		client:          client,
		priceRepository: priceRepository,
	}
}

func (p changeHeroProvider) ProvideAddressForTransfer(input provider.ProvideAddressForTransferInput) (*provider.ProvideAddressForTransferOutput, *app_error.AppError) {
	limits, limitsErr := p.GetLimits(provider.GetLimitsInput{
		CoinFromName:    input.CoinFromName,
		CoinFromNetwork: input.CoinFromNetwork,
		CoinToName:      input.CoinToName,
		CoinToNetwork:   input.CoinToNetwork,
	})
	if limitsErr != nil {
		return nil, limitsErr
	}

	if input.Amount < limits.Min || input.Amount > limits.Max {
		return nil, app_error.InvalidDataError(errors.Errorf("amount out of range. min: %f, max: %f", limits.Min, limits.Max))
	}

	from, err := p.transformCoinName(input.CoinFromName)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.Errorf("coin from `%s` does not supports exchange", input.CoinFromName))
	}

	to, err := p.transformCoinName(input.CoinToName)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.Errorf("coin to `%s` does not supports exchange", input.CoinToName))
	}

	response, err := p.client.CreateTransaction(change_hero.CreateTransactionParams{
		From:      from,
		To:        to,
		AddressTo: input.AddressTo,
		Amount:    input.Amount,
	})
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to create transaction: %s", err))
	}

	return &provider.ProvideAddressForTransferOutput{
		PayInAddress:  response.Result.PayInAddress,
		TransactionId: response.Result.Id,
		SupportLink:   "support@changehero.io",
	}, nil
}

func (p changeHeroProvider) GetLimits(input provider.GetLimitsInput) (*provider.GetLimitsOutput, *app_error.AppError) {
	from, err := p.transformCoinName(input.CoinFromName)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.Errorf("coin from `%s` does not supports exchange", input.CoinFromName))
	}

	to, err := p.transformCoinName(input.CoinToName)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.Errorf("coin to `%s` does not supports exchange", input.CoinToName))
	}

	response, err := p.client.GetMinAmount(change_hero.GetMinAmountParams{
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to get min amount: %s", err))
	}

	minAmount, err := strconv.ParseFloat(response.Result, 64)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to parse min amount: %s", err))
	}

	priceInUsd, err := p.priceRepository.FindLatest(coin.FindOptions{
		CoinName:     input.CoinFromName,
		CoinNetwork:  input.CoinFromNetwork,
		FiatCurrency: app_enum.UsdCurrency,
	}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to find latest price: %s", err))
	}

	floatPriceInUsd, err := app_util.AmountToFloat(input.CoinFromNetwork, priceInUsd.Price)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not format price to float: %s", err))
	}

	return &provider.GetLimitsOutput{
		Min: minAmount,
		Max: 1000 / floatPriceInUsd,
	}, nil
}

func (p changeHeroProvider) GetExchangeAmount(input provider.GetExchangeAmountInput) (*provider.GetExchangeAmountOutput, *app_error.AppError) {
	limits, limitsErr := p.GetLimits(provider.GetLimitsInput{
		CoinFromName:    input.CoinFromName,
		CoinFromNetwork: input.CoinFromNetwork,
		CoinToName:      input.CoinToName,
		CoinToNetwork:   input.CoinToNetwork,
	})
	if limitsErr != nil {
		return nil, limitsErr
	}

	if input.SendAmount < limits.Min || input.SendAmount > limits.Max {
		return nil, app_error.InvalidDataError(errors.Errorf("amount out of range. min: %f, max: %f", limits.Min, limits.Max))
	}

	from, err := p.transformCoinName(input.CoinFromName)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.Errorf("coin from `%s` does not supports exchange", input.CoinFromName))
	}

	to, err := p.transformCoinName(input.CoinToName)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.Errorf("coin to `%s` does not supports exchange", input.CoinToName))
	}

	response, err := p.client.GetExchangeAmount(change_hero.GetExchangeAmountParams{
		From:   from,
		To:     to,
		Amount: input.SendAmount,
	})
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to get exchange amount: %s", err))
	}

	receiveAmount, err := strconv.ParseFloat(response.Result, 64)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to parse exchange amount: %s", err))
	}

	return &provider.GetExchangeAmountOutput{ReceiveAmount: receiveAmount}, nil
}

func (p changeHeroProvider) transformCoinName(coinName app_enum.CoinName) (string, error) {
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
