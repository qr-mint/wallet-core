package dex

import (
	"encoding/base64"
	"fmt"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type PoolType int64

const (
	VOLATILE PoolType = 0 // TON
	STABLE   PoolType = 1
)

type PoolAssets struct {
	PoolType PoolType
	From     *address.Address
	To       *address.Address
}

func (pool *PoolAssets) FromAsset() *Asset {
	if pool.From != nil {
		return JettonAsset(pool.From)
	}
	return NativeAsset()
}

func (pool *PoolAssets) ToAsset() *Asset {
	if pool.To != nil {
		return JettonAsset(pool.To)
	}
	return NativeAsset()
}

func (pool *PoolAssets) Address() (*address.Address, error) {
	cell1, err := pool.FromAsset().toSlice().ToCell()
	if err != nil {
		return nil, err
	}
	cell2, err := pool.ToAsset().toSlice().ToCell()
	if err != nil {
		return nil, err
	}
	//

	stackItem := StackItem{
		Type:  "num",
		Value: fmt.Sprintf("0x%x", pool.PoolType),
	}
	stackItemTonOrJetton := StackItem{
		Type:  "slice",
		Value: base64.StdEncoding.EncodeToString(cell1.ToBOCWithFlags(false)),
	}
	stackItemOnlyJetton := StackItem{
		Type:  "slice",
		Value: base64.StdEncoding.EncodeToString(cell2.ToBOCWithFlags(false)),
	}
	data := DataStructure{
		Address: MAINNET_FACTORY_ADDR,
		Method:  "get_pool_address",
		Stack:   []StackItem{stackItem, stackItemTonOrJetton, stackItemOnlyJetton},
	}
	//
	result, err := RunGetMethod(data)
	if err != nil {
		return nil, err
	}
	//
	bytes, err := base64.StdEncoding.DecodeString(result.Stack[0].Value)
	if err != nil {
		return nil, err
	}

	addrCell, err := cell.FromBOC(bytes)
	if err != nil {
		return nil, err
	}
	poolAddress, err := addrCell.BeginParse().LoadAddr()
	if err != nil {
		return nil, fmt.Errorf("failed to load address from result slice: %w", err)
	}

	return poolAddress, nil
}
