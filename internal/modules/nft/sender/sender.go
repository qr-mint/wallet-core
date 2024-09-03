package sender

import "nexus-wallet/internal/app_error"

type NftSender interface {
	Send(message interface{}) (*SendOutput, *app_error.AppError)
}

type SendOutput struct {
	Hash string
}
