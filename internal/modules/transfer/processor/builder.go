package processor

import (
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
)

type TransferMessageBuilder interface {
	BuildMessage(input BuildMessageInput) (*BuildMessageOutput, *app_error.AppError)
	BuildTokenMessage(input BuildTokenMessageInput) (*BuildMessageOutput, *app_error.AppError)
}

type BuildMessageInput struct {
	FromAddress string
	ToAddress   string
	Amount      float64
	Options     BuildMessageWalletInputOptions
}

type BuildTokenMessageInput struct {
	FromAddress     string
	ToAddress       string
	Amount          float64
	ContractAddress string
	CoinName        app_enum.CoinName
	Options         BuildMessageWalletInputOptions
}

type BuildMessageWalletInputOptions struct {
	Version   uint8
	Comment   string
	Memo      string
	PublicKey string
}

type BuildMessageOutput struct {
	Message interface{}
}
