package sender

import "nexus-wallet/internal/app_error"

type NftMessageBuilder interface {
	BuildSendNftMessage(input BuildNftMessageInput) (*BuildMessageOutput, *app_error.AppError)
}

type BuildNftMessageInput struct {
	FromAddress string
	ToAddress   string
	NftAddress  string
	Options     BuildMessageWalletInputOptions
}

type BuildMessageWalletInputOptions struct {
	Version   uint8
	PublicKey string
}

type BuildMessageOutput struct {
	Message interface{}
}
