package dex

import (
	"encoding/base64"

	"fmt"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type Vault struct {
	Asset *Asset
}

func (state *Vault) GetVaultAddress() (*address.Address, error) {
	//var stack tlb.Stack
	assetCell, err := state.Asset.toSlice().ToCell()
	if err != nil {
		return nil, err
	}
	base64boc := base64.StdEncoding.EncodeToString(assetCell.ToBOCWithFlags(false))
	stackItem := StackItem{
		Type:  "slice",
		Value: base64boc,
	}
	data := DataStructure{
		Address: MAINNET_FACTORY_ADDR,
		Method:  "get_vault_address",
		Stack:   []StackItem{stackItem},
	}
	result, err := RunGetMethod(data)
	if err != nil {
		return nil, err
	}
	// Получаем адрес из стека
	bytes, err := base64.StdEncoding.DecodeString(result.Stack[0].Value)
	if err != nil {
		return nil, err
	}
	addrCell, err := cell.FromBOC(bytes)
	if err != nil {
		return nil, err
	}
	vaultAddress, err := addrCell.BeginParse().LoadAddr()
	if err != nil {
		return nil, fmt.Errorf("failed to load address from result slice: %w", err)
	}

	return vaultAddress, nil
}
