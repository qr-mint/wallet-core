package ton

import (
	"fmt"
	"gitlab.com/golib4/logger/logger"
	"gitlab.com/golib4/toncenter-client/toncenter"
	"nexus-wallet/internal/modules/transaction/provider"
	"time"
)

type toncenterProvideProvider struct {
	client        *toncenter.Client
	dataExtractor *ToncenterDataExtractor
	logger        logger.Logger
}

func NewToncenterProvider(
	client *toncenter.Client,
	dataExtractor *ToncenterDataExtractor,
	logger logger.Logger,
) provider.Provider {
	return &toncenterProvideProvider{
		client:        client,
		dataExtractor: dataExtractor,
		logger:        logger,
	}
}

func (p toncenterProvideProvider) GetTransactions(account string, limit int32, timestampFrom time.Time) ([]provider.BlockchainTransactionData, error) {
	var minTimestamp int64 = 0
	if timestampFrom.Unix() >= 0 {
		minTimestamp = timestampFrom.Unix()
	}

	transactions, err := p.client.GetAccountTransactions(toncenter.TransactionsRequest{
		Account:      account,
		Limit:        limit,
		MinTimestamp: minTimestamp,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet transactions: %s", err)
	}

	var result []provider.BlockchainTransactionData
	for _, transaction := range transactions.Transactions {
		extractedData, err := p.dataExtractor.Extract(transaction, account)
		if err != nil {
			p.logger.Warningf("failed to extract toncenter transaction data: %s", err)
			continue
		}
		if extractedData == nil {
			continue
		}

		result = append(result, *extractedData)
	}

	return result, nil
}
