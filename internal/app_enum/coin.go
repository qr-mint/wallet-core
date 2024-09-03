package app_enum

import (
	"nexus-wallet/internal/app_enum/utils"
)

type CoinName string

const (
	TonCoinName    CoinName = "ton"
	TetherCoinName CoinName = "tether"
	TronCoinName   CoinName = "tron"
)

func ToCoinName(value string) *CoinName {
	if !utils.AssertInArray(value, []string{string(TonCoinName), string(TetherCoinName), string(TronCoinName)}) {
		return nil
	}

	coinName := CoinName(value)
	return &coinName
}

func GetCoinNamesByNetwork(network Network) []CoinName {
	switch network {
	case Trc20Network:
		return []CoinName{TronCoinName, TetherCoinName}
	case TonNetwork:
		return []CoinName{TonCoinName}
	}

	return nil
}
