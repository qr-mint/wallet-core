package dex

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	TRANSFER = 0xf8a7ea5
)

type JettonBody struct {
	QueryId         uint64
	Amount          tlb.Coins
	Destination     *address.Address
	ResponseAddress *address.Address
	CustomPayload   *cell.Cell
	ForwardAmount   tlb.Coins
	ForwardPayload  *cell.Cell
}

func (payload *JettonBody) ToCell() *cell.Cell {
	return cell.BeginCell().
		MustStoreUInt(TRANSFER, 32).
		MustStoreUInt(payload.QueryId, 64).
		MustStoreCoins(payload.Amount.Nano().Uint64()).
		MustStoreAddr(payload.Destination).
		MustStoreAddr(payload.ResponseAddress).
		MustStoreMaybeRef(payload.CustomPayload).
		MustStoreCoins(payload.ForwardAmount.Nano().Uint64()).
		MustStoreRef(payload.ForwardPayload).
		EndCell()
}
