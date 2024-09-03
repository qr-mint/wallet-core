package trc20

import (
	"fmt"
	"gitlab.com/golib4/logger/logger"
	"gitlab.com/golib4/trongrid-client/trongrid"
	"gitlab.com/golib4/tronscanapi-client/tronscanapi"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/transaction/provider"
	"sync"
	"time"
)

type tronProvider struct {
	trongridClient    *trongrid.Client
	tronscanapiClient *tronscanapi.Client
	dataExtractor     *TronDataExtractor
	logger            logger.Logger
}

func NewTronProvider(
	client *trongrid.Client,
	dataExtractor *TronDataExtractor,
	tronscanapiClient *tronscanapi.Client,
	logger logger.Logger,
) provider.Provider {
	return &tronProvider{
		trongridClient:    client,
		dataExtractor:     dataExtractor,
		tronscanapiClient: tronscanapiClient,
		logger:            logger,
	}
}

func (p tronProvider) GetTransactions(account string, limit int32, timestampFrom time.Time) ([]provider.BlockchainTransactionData, error) {
	var wg sync.WaitGroup
	results := make(chan provider.BlockchainTransactionData)
	errors := make(chan error, 2)

	var minTimestamp int64 = 0
	if timestampFrom.Unix() >= 0 {
		minTimestamp = timestampFrom.Unix()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.fetchTokenTransactions(account, limit, minTimestamp, results, errors)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.fetchAccountTransactions(account, limit, minTimestamp, results, errors)
	}()

	go func() {
		wg.Wait()
		defer close(results)
		defer close(errors)
	}()

	select {
	case err := <-errors:
		return nil, fmt.Errorf("error fetching transactions: %s", err)
	default:
		var finalResult []provider.BlockchainTransactionData
		for result := range results {
			finalResult = append(finalResult, result)
		}

		return finalResult, nil
	}
}

func (p tronProvider) fetchAccountTransactions(account string, limit int32, minTimestamp int64, results chan<- provider.BlockchainTransactionData, errors chan<- error) {
	transactions, err := p.trongridClient.GetWalletTransactions(trongrid.GetWalletTransactionRequest{
		WalletAddress: account,
		Query:         trongrid.GetWalletTransactionRequestQuery{Limit: limit, MitTimestamp: time.Unix(minTimestamp, 0).UnixMilli()},
	})

	if err != nil {
		errors <- fmt.Errorf("failed to get wallet transactions: %s", err)
		return
	}

	for _, transaction := range transactions.Data {
		extractedData, err := p.dataExtractor.Extract(transaction, account)
		if err != nil {
			p.logger.Warningf("failed to extract token transaction data: %s", err)
			continue
		}
		if extractedData == nil || extractedData.CoinName != app_enum.TronCoinName {
			continue
		}

		results <- *extractedData
	}
}

func (p tronProvider) fetchTokenTransactions(account string, limit int32, minTimestamp int64, results chan<- provider.BlockchainTransactionData, errors chan<- error) {
	transactions, err := p.tronscanapiClient.GetTokenTransfers(tronscanapi.GetTokenTransfersRequests{
		Limit:          limit,
		StartTimestamp: int32(time.Unix(minTimestamp, 0).UnixMilli()),
		Address:        account,
	})
	if err != nil {
		errors <- fmt.Errorf("failed to get wallet token transactions: %s", err)
		return
	}

	for _, transaction := range transactions.TokenTransfers {
		extractedData, err := p.dataExtractor.ExtractToken(transaction, account)
		if err != nil {
			p.logger.Warningf("failed to extract tron token transaction data: %s", err)
			continue
		}
		if extractedData == nil {
			continue
		}

		results <- *extractedData
	}
}
