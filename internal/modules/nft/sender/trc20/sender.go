package trc20

import (
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	sender_module "nexus-wallet/internal/modules/nft/sender"
)

type sender struct {
}

func NewSender() sender_module.NftSender {
	return sender{}
}

func (sender) Send(message interface{}) (*sender_module.SendOutput, *app_error.AppError) {
	return nil, app_error.IllegalOperationError(errors.New("trc20 Send nft is not implemented"))
}
