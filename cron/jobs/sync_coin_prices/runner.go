package sync_coin_prices

import (
	"nexus-wallet/internal/modules/coin"
)

type Runner struct {
	coinService *coin.Service
}

func NewRunner(coinService *coin.Service) *Runner {
	return &Runner{
		coinService: coinService,
	}
}

func (r Runner) Run() {
	r.coinService.SyncPrices()
}

func (Runner) GetPattern() string {
	return "0 */3 * * *" // every 3 hours
}
