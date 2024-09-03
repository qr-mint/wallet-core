package app_util

import (
	"fmt"
	tron_address "github.com/fbsobreira/gotron-sdk/pkg/address"
	ton_address "github.com/xssnick/tonutils-go/address"
	"nexus-wallet/internal/app_enum"
)

func IsValidAddress(addressString string, network app_enum.Network) (bool, error) {
	switch network {
	case app_enum.TonNetwork:
		if _, err := ton_address.ParseAddr(addressString); err != nil {
			return false, nil
		}
		return true, nil
	case app_enum.Trc20Network:
		if _, err := tron_address.Base58ToAddress(addressString); err != nil {
			return false, nil
		}
		return true, nil
	}

	return false, fmt.Errorf("invalid network provided: %s", network)
}
