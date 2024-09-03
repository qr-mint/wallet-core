package provider

import "nexus-wallet/internal/app_error"

type NftProvider interface {
	Provide(input ProvideInput) (*ProvideOutput, *app_error.AppError)
}

type ProvideInput struct {
	OwnerAddress string
}

type NftItem struct {
	Price                 int64
	PriceTokenSymbol      string
	PreviewsData          []byte
	Index                 int64
	Address               string
	Name                  string
	CollectionName        string
	CollectionAddress     string
	CollectionDescription string
}

type ProvideOutput struct {
	Items []NftItem
}
