package trc20

import (
	"fmt"
	client "github.com/fbsobreira/gotron-sdk/pkg/client"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/address/balance/provider"
	"strings"
)

type trongridGrpcProvider struct {
	grpcClient *client.GrpcClient
}

func NewTrongridGrpcProvider(grpcClient *client.GrpcClient) provider.Provider {
	return &trongridGrpcProvider{
		grpcClient: grpcClient,
	}
}

func (p trongridGrpcProvider) GetMainCoinBalance(accountAddress string) (*provider.GetBalanceOutput, *app_error.AppError) {
	response, err := p.grpcClient.GetAccount(accountAddress)
	if err != nil {
		if strings.Contains(err.Error(), "account not found") {
			return &provider.GetBalanceOutput{Value: 0}, nil
		}
		return nil, app_error.InternalError(fmt.Errorf("can not get balance in trongridGrpcProvider: %s", err))
	}

	return &provider.GetBalanceOutput{Value: response.Balance}, nil
}

func (p trongridGrpcProvider) GetCoinBalance(accountAddress string, tokenAddress string) (*provider.GetBalanceOutput, *app_error.AppError) {
	balanceAmount, err := p.grpcClient.TRC20ContractBalance(accountAddress, tokenAddress)
	if err != nil {
		if strings.Contains(err.Error(), "account not found") {
			return &provider.GetBalanceOutput{Value: 0}, nil
		}
		return nil, app_error.InternalError(fmt.Errorf("can not get token in trongridGrpcProvider: %s", err))
	}
	return &provider.GetBalanceOutput{Value: balanceAmount.Int64()}, nil
}
