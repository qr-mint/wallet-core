package provider

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
)

type TransferableProvider interface {
	ProvideAddressForTransfer(input ProvideAddressForTransferInput) (*ProvideAddressForTransferOutput, *app_error.AppError)
	GetLimits(input GetLimitsInput) (*GetLimitsOutput, *app_error.AppError)
	GetExchangeAmount(input GetExchangeAmountInput) (*GetExchangeAmountOutput, *app_error.AppError)
}

type ProvideAddressForTransferInput struct {
	CoinFromName    app_enum.CoinName
	CoinFromNetwork app_enum.Network

	CoinToName    app_enum.CoinName
	CoinToNetwork app_enum.Network

	AddressTo string
	Amount    float64
}

type ProvideAddressForTransferOutput struct {
	PayInAddress  string
	TransactionId string
	SupportLink   string
}

type GetLimitsInput struct {
	CoinFromName    app_enum.CoinName
	CoinFromNetwork app_enum.Network

	CoinToName    app_enum.CoinName
	CoinToNetwork app_enum.Network
}

type GetLimitsOutput struct {
	Min float64
	Max float64
}

type GetExchangeAmountInput struct {
	CoinFromName    app_enum.CoinName
	CoinFromNetwork app_enum.Network

	CoinToName    app_enum.CoinName
	CoinToNetwork app_enum.Network

	SendAmount float64
}

type GetExchangeAmountOutput struct {
	ReceiveAmount float64
}
