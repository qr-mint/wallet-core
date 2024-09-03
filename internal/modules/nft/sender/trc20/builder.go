package trc20

import (
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	sender_module "nexus-wallet/internal/modules/nft/sender"
)

type builder struct {
}

func NewBuilder() sender_module.NftMessageBuilder {
	return builder{}
}

func (builder) BuildSendNftMessage(input sender_module.BuildNftMessageInput) (*sender_module.BuildMessageOutput, *app_error.AppError) {
	return nil, app_error.IllegalOperationError(errors.New("trc20 BuildSendNftMessage is not implemented"))
}
