package address

import "nexus-wallet/internal/app_enum"

type AddressCoin struct {
	AddressCoinId int64
	WalletAddress string
	CoinName      app_enum.CoinName
	CoinNetwork   app_enum.Network
}
