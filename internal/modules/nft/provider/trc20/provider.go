package trc20

import (
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	provider_module "nexus-wallet/internal/modules/nft/provider"
)

type provider struct {
}

func NewProvider() provider_module.NftProvider {
	return &provider{}
}

func (p provider) Provide(input provider_module.ProvideInput) (*provider_module.ProvideOutput, *app_error.AppError) {
	return nil, app_error.IllegalOperationError(errors.New("trc20 nft Provide is not implemented"))
}
