package provider

import "nexus-wallet/internal/app_error"

type Provider interface {
	GetMainCoinBalance(accountAddress string) (*GetBalanceOutput, *app_error.AppError)
	GetCoinBalance(accountAddress string, tokenAddress string) (*GetBalanceOutput, *app_error.AppError)
}

type GetBalanceOutput struct {
	Value int64
}
