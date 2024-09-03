package util

import (
	"fmt"
	"nexus-wallet/internal/app_enum"
)

func ProvideExplorerLink(network app_enum.Network, hash string) (string, error) {
	switch network {
	case app_enum.TonNetwork:
		return fmt.Sprintf("https://tonscan.org/tx/%s", hash), nil
	case app_enum.Trc20Network:
		return fmt.Sprintf("https://tronscan.org/#/transaction/%s", hash), nil
	default:
		return "", fmt.Errorf("can not resolve explorer link. invalid network provided: %s", network)
	}
}
