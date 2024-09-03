package address

import "nexus-wallet/internal/app_enum"

type AddressCoin struct {
	WalletAddress string
	CoinName      app_enum.CoinName
}
