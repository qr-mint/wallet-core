package processor

import "nexus-wallet/internal/app_error"

type TransferProcessor interface {
	Transfer(message interface{}) (*TransferOutput, *app_error.AppError)
}

type TransferOutput struct {
	Hash string
}
