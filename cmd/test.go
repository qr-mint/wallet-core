package main

import (
	"log"
	"nexus-wallet/cmd/dex"

	"github.com/xssnick/tonutils-go/address"
)

func main() {
	SCALE_ADDRESS := address.MustParseAddr("EQBlqsm144Dq6SjbPI4jjZvA1hqTIP3CvHovbIfW_t-SCALE")

	TONAsset := dex.NativeAsset()
	// SCALE := dex.JettonAsset(SCALE_ADDRESS)

	pool := dex.PoolAssets{
		PoolType: dex.VOLATILE,
		To:       SCALE_ADDRESS,
	}
	addr, err := pool.Address()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(addr.String()) // https://tonscan.org/jetton/EQDcm06RlreuMurm-yik9WbL6kI617B77OrSRF_ZjoCYFuny

	vault := dex.Vault{
		Asset: TONAsset,
	}
	vaultAddress, err := vault.GetVaultAddress()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(addr.String())
}
