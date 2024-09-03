package app_enum

import (
	"nexus-wallet/internal/app_enum/utils"
)

type Network string

const (
	TonNetwork   Network = "ton"
	Trc20Network Network = "trc20"
)

func ToNetwork(value string) *Network {
	if !utils.AssertInArray(value, []string{string(TonNetwork), string(Trc20Network)}) {
		return nil
	}

	network := Network(value)
	return &network
}

func GetNetworks() []Network {
	return []Network{
		TonNetwork,
		Trc20Network,
	}
}
