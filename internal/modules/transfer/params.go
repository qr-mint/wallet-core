package transfer

import "nexus-wallet/internal/app_enum"

type BuildTransferInput struct {
	Network    app_enum.Network
	CoinName   app_enum.CoinName
	MnemonicId int64
	ToAddress  string
	Amount     float64
	PublicKey  string
	Version    uint8
	Comment    string
}

type BuildTransferOutput struct {
	Message interface{}
}

type TransferInput struct {
	Network    app_enum.Network
	CoinName   app_enum.CoinName
	Amount     float64
	MnemonicId int64
	Message    interface{}
}

type TransferOutput struct {
	Hash string
}
