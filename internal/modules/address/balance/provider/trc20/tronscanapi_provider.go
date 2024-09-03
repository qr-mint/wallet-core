package trc20

import (
	"fmt"
	"gitlab.com/golib4/tronscanapi-client/tronscanapi"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/address/balance/provider"
	"strconv"
)

type tronscanapiProvider struct {
	client *tronscanapi.Client
}

func NewTronscanapiProvider(client *tronscanapi.Client) provider.Provider {
	return &tronscanapiProvider{
		client: client,
	}
}
func (p tronscanapiProvider) GetMainCoinBalance(accountAddress string) (*provider.GetBalanceOutput, *app_error.AppError) {
	response, err := p.client.GetAddressBalance(tronscanapi.GetAddressBalanceRequest{Address: accountAddress})
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get balance in tronscanapiProvider: %s", err))
	}

	for _, token := range response.Data {
		if token.TokenName == "trx" {
			balance, err := strconv.Atoi(token.Amount)
			if err != nil {
				return nil, app_error.InternalError(fmt.Errorf("can not parce balance in tronscanapiProvider: %s", err))
			}

			return &provider.GetBalanceOutput{Value: int64(balance)}, nil
		}
	}

	return &provider.GetBalanceOutput{Value: 0}, nil
}

func (p tronscanapiProvider) GetCoinBalance(accountAddress string, tokenAddress string) (*provider.GetBalanceOutput, *app_error.AppError) {
	response, err := p.client.GetAddressBalance(tronscanapi.GetAddressBalanceRequest{Address: accountAddress})
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get balance in tronscanapiProvider: %s", err))
	}

	for _, token := range response.Data {
		if token.TokenName == "Tether USD" {
			balance, err := strconv.Atoi(token.Amount)
			if err != nil {
				return nil, app_error.InternalError(fmt.Errorf("can not parce coin balance in tronscanapiProvider: %s", err))
			}

			return &provider.GetBalanceOutput{Value: int64(balance)}, nil
		}
	}

	return &provider.GetBalanceOutput{Value: 0}, nil
}
