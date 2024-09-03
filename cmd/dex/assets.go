package dex

import (
	"fmt"

	"github.com/xssnick/tonutils-go/address"

	"github.com/xssnick/tonutils-go/tvm/cell"
)

type AssetType uint8

const (
	Native AssetType = 0b0000 // TON
	Jetton AssetType = 0b0001
)

type Asset struct {
	Type    AssetType
	Address *address.Address
}

func NewAsset(assetType AssetType, address *address.Address) *Asset {
	return &Asset{
		Type:    assetType,
		Address: address,
	}
}

func NativeAsset() *Asset {
	return NewAsset(Native, nil)
}

func JettonAsset(minter *address.Address) *Asset {
	return NewAsset(Jetton, minter)
}

func (a *Asset) Equals(other *Asset) bool {
	return a.ToString() == other.ToString()
}

func (a *Asset) ToString() string {
	switch a.Type {
	case Native:
		return "native"

	case Jetton:
		if a.Address != nil {
			return fmt.Sprintf("jetton:%s", a.Address.String())
		}
	}

	return ""
}

func (a *Asset) toSlice() *cell.Slice {
	var builder *cell.Builder
	if a.Type == Native {
		builder = new(cell.Builder).MustStoreUInt(uint64(a.Type), 4)
	} else if a.Type == Jetton {
		builder = new(cell.Builder).
			MustStoreUInt(uint64(a.Type), 4).
			MustStoreInt(int64(a.Address.Workchain()), 8).
			MustStoreBinarySnake(a.Address.Data())

	}
	return cell.BeginCell().MustStoreBuilder(builder).EndCell().BeginParse()
	//beginCell().storeWritable(this).endCell().beginParse();
}
