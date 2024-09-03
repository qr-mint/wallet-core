package dex

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"gitlab.com/golib4/coins/coins"
)

const DEPOSIT_LIQUIDITY = 0xd55e4686
const SWAP = 0xea06185d
const CREATE_VAULT = 0x21cfe02b
const CREATE_VOLATILE_POOL = 0x97d51f2f

type SwapStep struct {
	PoolAddress *address.Address
	Limit       uint64
	Next        *SwapStep
}

func (swapStep *SwapStep) ToCell() *cell.Cell {
	return cell.BeginCell().
		MustStoreAddr(swapStep.PoolAddress).
		MustStoreUInt(0, 1).
		MustStoreCoins(swapStep.Limit).
		MustStoreMaybeRef(swapStep.NextToCell()).
		EndCell()
}

func (swapStep *SwapStep) NextToCell() *cell.Cell {
	if swapStep.Next != nil {
		return swapStep.Next.ToCell()
	}
	return nil
}

type SwapParams struct {
	Deadline         uint64
	RecipientAddress *address.Address
	ReferralAddress  *address.Address
	FulfillPayload   *cell.Cell
	RejectPayload    *cell.Cell
}

func (swapParams *SwapParams) ToCell() *cell.Cell {
	beginCell := cell.BeginCell()
	if swapParams.Deadline != 0 {
		beginCell.MustStoreUInt(swapParams.Deadline, 0)
	} else {
		beginCell.MustStoreUInt(0, 32)
	}

	return beginCell.MustStoreAddr(swapParams.RecipientAddress).
		MustStoreAddr(swapParams.ReferralAddress).
		MustStoreMaybeRef(swapParams.FulfillPayload).
		MustStoreMaybeRef(swapParams.RejectPayload).
		EndCell()
}

type SwapBody struct {
	QueryId     uint64
	Amount      *coins.Coins
	PoolAddress *address.Address
	Limit       *coins.Coins
	SwapParams  SwapParams
	Next        SwapStep
	GasAmount   uint64
}

func (input *SwapBody) ToCell() *cell.Cell {
	return cell.BeginCell().
		MustStoreUInt(SWAP, 32).
		MustStoreUInt(input.QueryId, 64).
		MustStoreCoins(input.Amount.Nano().Uint64()).
		MustStoreAddr(input.PoolAddress).
		MustStoreUInt(0, 1).
		MustStoreCoins(input.Limit.Nano().Uint64()).
		MustStoreMaybeRef(input.Next.ToCell()).
		MustStoreRef(input.SwapParams.ToCell()).
		EndCell()
}

type CreateVaultBody struct {
	QueryId uint64
	Asset   Asset
}

func (body *CreateVaultBody) ToCell() *cell.Cell {
	beginCell := cell.BeginCell().
		MustStoreUInt(CREATE_VAULT, 32)
	if body.QueryId != 0 {
		beginCell.MustStoreUInt(body.QueryId, 0)
	} else {
		beginCell.MustStoreUInt(0, 64)
	}
	slice := body.Asset.toSlice()
	bytes := slice.MustLoadBinarySnake()
	return beginCell.MustStoreSlice(bytes, slice.BitsLeft()).EndCell()
}

type CreateVaultPoolBody struct {
	QueryId uint64
	Asset   []Asset
}

func (body *CreateVaultPoolBody) ToCell() *cell.Cell {
	beginCell := cell.BeginCell().
		MustStoreUInt(CREATE_VOLATILE_POOL, 32)
	if body.QueryId != 0 {
		beginCell.MustStoreUInt(body.QueryId, 0)
	} else {
		beginCell.MustStoreUInt(0, 64)
	}
	slice1 := body.Asset[0].toSlice()
	bytes1 := slice1.MustLoadBinarySnake()
	slice2 := body.Asset[0].toSlice()
	bytes2 := slice1.MustLoadBinarySnake()
	return beginCell.
		MustStoreSlice(bytes1, slice1.BitsLeft()).
		MustStoreSlice(bytes2, slice2.BitsLeft()).
		EndCell()
}

// type JettonStepBody struct {
// 	QueryId uint64
//   Destination *address.Address
//     FiatAmount uint64
//   ResponseAddress *address.Address
//       CustomPayload *cell.Cell
//       ForwardAmount uint64
//     ForwardPayload *cell.Cell
// }

// func (body *JettonStepBody) ToCell() *cell.Cell {
//  return cell.BeginCell().
//         storeUint(JettonWallet.TRANSFER, 32).
//         storeUint(queryId ?? 0, 64).
//         storeCoins(amount).
//         storeAddress(destination).
//         storeAddress(responseAddress).
//         storeMaybeRef(customPayload).
//         storeCoins(forwardAmount ?? 0).
//         storeMaybeRef(forwardPayload)
//         endCell();
// }
