package ton

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xssnick/tonutils-go/address"
	"gitlab.com/golib4/tonconsole-client/tonconsole"
	"nexus-wallet/internal/app_error"
	provider_module "nexus-wallet/internal/modules/nft/provider"
	"strconv"
)

type provider struct {
	client *tonconsole.Client
}

func NewProvider(client *tonconsole.Client) provider_module.NftProvider {
	return &provider{client: client}
}

func (p provider) Provide(input provider_module.ProvideInput) (*provider_module.ProvideOutput, *app_error.AppError) {
	data, err := p.client.GetAccountNft(input.OwnerAddress)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get wallet nfts list: %s", err))
	}

	nftItems, err := p.transformResponse(*data)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not transform response to NFT data: %s", err))
	}

	return &provider_module.ProvideOutput{Items: nftItems}, nil
}

func (p provider) transformResponse(response tonconsole.NftsResponse) ([]provider_module.NftItem, error) {
	var items []provider_module.NftItem
	for _, dataItem := range response.Items {
		var int64Price int64
		if dataItem.Sale.Price.Value != "" {
			intPrice, err := strconv.Atoi(dataItem.Sale.Price.Value)
			if err != nil {
				return nil, fmt.Errorf("can not convert price to int: %s", err)
			}
			int64Price = int64(intPrice)
		}
		var collectionAddress string
		if dataItem.Collection.Address != "" {
			collectionAddress = address.MustParseRawAddr(dataItem.Collection.Address).Bounce(false).String()
		}
		var collectionName string
		if dataItem.Collection.Name != "" {
			collectionName = dataItem.Collection.Name
		}
		var collectionDescription string
		if dataItem.Collection.Description != "" {
			collectionDescription = dataItem.Collection.Description
		}
		previews, err := json.Marshal(dataItem.Previews)
		if err != nil {
			return nil, errors.Errorf("can not marshal previews: %s", err)
		}
		items = append(items, provider_module.NftItem{
			Price:                 int64Price,
			PriceTokenSymbol:      dataItem.Sale.Price.TokenName,
			PreviewsData:          previews,
			Index:                 dataItem.Index,
			Address:               address.MustParseRawAddr(dataItem.Address).Bounce(false).String(),
			Name:                  dataItem.Metadata.Name,
			CollectionName:        collectionName,
			CollectionAddress:     collectionAddress,
			CollectionDescription: collectionDescription,
		})
	}

	return items, nil
}
