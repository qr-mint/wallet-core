package app_util

import (
	"fmt"
	"gitlab.com/golib4/coins/coins"
	"nexus-wallet/internal/app_enum"
)

func AmountToFloat(network app_enum.Network, amount int64) (float64, error) {
	switch network {
	case app_enum.TonNetwork:
		result, _ := coins.FromNanoTokenU(uint64(amount), 9).Normolized().Float64()
		return result, nil
	case app_enum.Trc20Network:
		result, _ := coins.FromNanoTokenU(uint64(amount), 6).Normolized().Float64()
		return result, nil
	default:
		return 0, fmt.Errorf("invalid network provided: %s", network)
	}
}

func AmountToInt(network app_enum.Network, amount float64) (int64, error) {
	switch network {
	case app_enum.TonNetwork:
		return int64(amount * 1000000000), nil
	case app_enum.Trc20Network:
		return int64(amount * 1000000), nil
	default:
		return 0, fmt.Errorf("invalid network provided: %s", network)
	}
}
