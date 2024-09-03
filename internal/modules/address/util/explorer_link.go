package util

import (
	"fmt"
	"nexus-wallet/internal/app_enum"
)

func ProvideExplorerLink(network app_enum.Network, address string) (string, error) {
	switch network {
	case app_enum.TonNetwork:
		return fmt.Sprintf("https://tonscan.org/address/%s", address), nil
	case app_enum.Trc20Network:
		return fmt.Sprintf("https://tronscan.org/#/address/%s", address), nil
	default:
		return "", fmt.Errorf("can not resolve explorer link. invalid network provided: %s", network)
	}
}
