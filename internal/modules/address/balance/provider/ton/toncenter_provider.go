package ton

import (
	"errors"
	"fmt"
	"gitlab.com/golib4/toncenter-client/toncenter"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/address/balance/provider"
	"strconv"
)

type toncenterProvider struct {
	client *toncenter.Client
}

func NewToncenterProvider(client *toncenter.Client) provider.Provider {
	return &toncenterProvider{
		client: client,
	}
}
func (p toncenterProvider) GetMainCoinBalance(accountAddress string) (*provider.GetBalanceOutput, *app_error.AppError) {
	response, err := p.client.GetAddressBalance(toncenter.GetAddressBalanceRequest{Address: accountAddress})
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get balance in ton provider: %s", err))
	}

	var balance int
	switch value := response.Result.(type) {
	case string:
		balance, err = strconv.Atoi(value)
		if err != nil {
			return nil, app_error.InternalError(fmt.Errorf("can not parce balance in ton provider: %s", err))
		}
	case int:
		balance = value
	default:
		return nil, app_error.InternalError(fmt.Errorf("can not parce balance in ton provider because of unexpected type: %T", response.Result))
	}

	return &provider.GetBalanceOutput{Value: int64(balance)}, nil
}

func (p toncenterProvider) GetCoinBalance(accountAddress string, tokenAddress string) (*provider.GetBalanceOutput, *app_error.AppError) {
	return nil, app_error.IllegalOperationError(errors.New("ton GetCoinBalance not implemented"))
}
