package nft

import (
	"encoding/json"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/nft/model/address"
	"nexus-wallet/internal/modules/nft/model/nft"
)

type ListInput struct {
	MnemonicId int64
	Limit      uint
	Offset     uint
}

type TonPreviewUrl struct {
	URL        string `json:"url"`
	Resolution string `json:"resolution"`
}

type ListOutputItem struct {
	Id                    int64
	Address               string
	Name                  string
	Price                 int64
	TokenSymbol           string
	Index                 int64
	CollectionAddress     string
	CollectionName        string
	CollectionDescription string
	Network               app_enum.Network
	PreviewUrls           interface{}
}

func (ListOutputItem) fillFromModel(nft nft.Nft, addressData address.Address) (*ListOutputItem, error) {
	outputItem := ListOutputItem{
		Id:                    nft.Id,
		Address:               nft.Address,
		Name:                  nft.Name,
		Price:                 nft.Price,
		TokenSymbol:           nft.TokenSymbol,
		Index:                 nft.Index,
		CollectionAddress:     nft.CollectionAddress,
		CollectionName:        nft.CollectionName,
		CollectionDescription: nft.CollectionDescription,
		Network:               addressData.Network,
	}

	switch addressData.Network {
	case app_enum.TonNetwork:
		var tonPreviewUrl []TonPreviewUrl
		err := json.Unmarshal(nft.PreviewData, &tonPreviewUrl)
		if err != nil {
			return nil, errors.Errorf("failed to unmarshal Ton PreviewUrl: %s", err)
		}
		outputItem.PreviewUrls = tonPreviewUrl
	default:
		return nil, errors.Errorf("unknown nft network: %s", addressData.Network)
	}

	return &outputItem, nil
}

type ListOutput struct {
	Items []ListOutputItem
}

func (ListOutput) fillFromModel(nftList []*nft.Nft, addresses map[int64]address.Address) (*ListOutput, error) {
	var outputItems []ListOutputItem
	for _, nftData := range nftList {
		outputItem, err := ListOutputItem{}.fillFromModel(*nftData, addresses[nftData.AddressId])
		if err != nil {
			return nil, errors.Errorf("failed to fill nft item: %s", err)
		}
		outputItems = append(outputItems, *outputItem)
	}

	return &ListOutput{Items: outputItems}, nil
}

type GetInput struct {
	MnemonicId int64
	NftId      int64
}

type GetOutput struct {
	Id                    int64
	Address               string
	Name                  string
	Price                 int64
	TokenSymbol           string
	Index                 int64
	CollectionAddress     string
	CollectionName        string
	CollectionDescription string
	Network               app_enum.Network
	PreviewUrls           interface{}
}

func (GetOutput) fillFromModel(nft nft.Nft, addressData address.Address) (*GetOutput, error) {
	outputItem := GetOutput{
		Id:                    nft.Id,
		Address:               nft.Address,
		Name:                  nft.Name,
		Price:                 nft.Price,
		TokenSymbol:           nft.TokenSymbol,
		Index:                 nft.Index,
		CollectionAddress:     nft.CollectionAddress,
		CollectionName:        nft.CollectionName,
		CollectionDescription: nft.CollectionDescription,
		Network:               addressData.Network,
		PreviewUrls:           nft.PreviewData,
	}

	switch addressData.Network {
	case app_enum.TonNetwork:
		var tonPreviewUrl []TonPreviewUrl
		err := json.Unmarshal(nft.PreviewData, &tonPreviewUrl)
		if err != nil {
			return nil, errors.Errorf("failed to unmarshal Ton PreviewUrl: %s", err)
		}
		outputItem.PreviewUrls = tonPreviewUrl
	default:
		return nil, errors.Errorf("unknown nft network: %s", addressData.Network)
	}

	return &outputItem, nil
}

type BuildSendMessageInput struct {
	Network    app_enum.Network
	MnemonicId int64
	ToAddress  string
	NftId      int64
	PublicKey  string
	Version    uint8
}

type BuildSendMessageOutput struct {
	Message interface{}
}

type SendInput struct {
	NftId      int64
	MnemonicId int64
	Network    app_enum.Network
	Message    interface{}
}

type SendOutput struct {
	Hash string
}
