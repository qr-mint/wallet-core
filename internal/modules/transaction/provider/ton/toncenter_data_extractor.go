package ton

import (
	"encoding/hex"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"gitlab.com/golib4/toncenter-client/toncenter"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/transaction/enum"
	"nexus-wallet/internal/modules/transaction/provider"
	"strconv"
	"strings"
	"time"
)

type ToncenterDataExtractor struct {
}

func NewToncenterDataExtractor() *ToncenterDataExtractor {
	return &ToncenterDataExtractor{}
}

func (e *ToncenterDataExtractor) Extract(transaction toncenter.Transaction, account string) (*provider.BlockchainTransactionData, error) {
	outcomingMessagesCount := len(transaction.OutMsgs)
	if outcomingMessagesCount == 0 && transaction.InMsg == nil {
		return nil, nil
	}
	if outcomingMessagesCount == 0 {
		return e.extractFromIncomingTransfer(transaction, account)
	}
	if outcomingMessagesCount == 1 && transaction.InMsg.Source == nil {
		if transaction.Account == transaction.InMsg.Destination {
			return e.extractFromOutcomingTransfer(transaction, account)
		}
	}

	return nil, fmt.Errorf("unsupportable transaction type. transaction hash: %s", transaction.Hash)

}

func (e *ToncenterDataExtractor) extractFromOutcomingTransfer(transaction toncenter.Transaction, account string) (*provider.BlockchainTransactionData, error) {
	sourceAddress, err := e.formatRawAddress(transaction.BlockRef.Workchain, transaction.Account)
	if err != nil {
		return nil, fmt.Errorf("can not resolve source address: %s", err)
	}
	destinationAddress, err := e.formatRawAddress(transaction.BlockRef.Workchain, transaction.OutMsgs[0].Destination)
	if err != nil {
		return nil, fmt.Errorf("can not resolve destination address: %s", err)
	}
	amountInt, err := strconv.Atoi(transaction.OutMsgs[0].Value)
	if err != nil {
		return nil, fmt.Errorf("can not convert out msgs amount to int: %s", err)
	}

	data := &provider.BlockchainTransactionData{
		From:      sourceAddress,
		To:        destinationAddress,
		Hash:      transaction.Hash,
		Amount:    int64(amountInt),
		CoinName:  app_enum.TonCoinName,
		Network:   app_enum.TonNetwork,
		CreatedAt: time.Unix(0, transaction.Now*int64(time.Second)),
	}
	e.setStatus(transaction, data)
	e.setType(account, data)

	return data, nil
}

func (e *ToncenterDataExtractor) extractFromIncomingTransfer(transaction toncenter.Transaction, account string) (*provider.BlockchainTransactionData, error) {
	sourceAddress := ""
	if transaction.InMsg.Source != nil {
		var err error
		sourceAddress, err = e.formatRawAddress(transaction.BlockRef.Workchain, *transaction.InMsg.Source)
		if err != nil {
			return nil, fmt.Errorf("can not resolve source address: %s", err)
		}
	}
	destinationAddress, err := e.formatRawAddress(transaction.BlockRef.Workchain, transaction.InMsg.Destination)
	if err != nil {
		return nil, fmt.Errorf("can not resolve destination address: %s", err)
	}

	var amountInt int
	if transaction.InMsg.Value != nil {
		amountInt, err = strconv.Atoi(*transaction.InMsg.Value)
		if err != nil {
			return nil, fmt.Errorf("can not convert in msg amount to int: %s", err)
		}
	}

	data := &provider.BlockchainTransactionData{
		From:      sourceAddress,
		To:        destinationAddress,
		Hash:      transaction.Hash,
		Amount:    int64(amountInt),
		CoinName:  app_enum.TonCoinName,
		Network:   app_enum.TonNetwork,
		CreatedAt: time.Unix(0, transaction.Now*int64(time.Second)),
	}
	e.setStatus(transaction, data)
	e.setType(account, data)

	return data, nil
}

func (e *ToncenterDataExtractor) setType(account string, data *provider.BlockchainTransactionData) {
	transactionType := enum.InType
	if data.From == address.MustParseAddr(account).Bounce(false).String() {
		transactionType = enum.OutType
	}

	data.Type = transactionType
}

func (e *ToncenterDataExtractor) setStatus(transaction toncenter.Transaction, data *provider.BlockchainTransactionData) {
	newWalletTxSuccess := (transaction.Description == nil || transaction.Description.Action == nil || transaction.Description.Action.ResultCode == nil) &&
		(transaction.Description == nil || transaction.Description.ComputePH == nil || transaction.Description.ComputePH.ExitCode == nil)

	executionSuccess := transaction.Description != nil &&
		transaction.Description.Action != nil &&
		transaction.Description.Action.ResultCode != nil &&
		*transaction.Description.Action.ResultCode <= 1

	isSuccess := newWalletTxSuccess || executionSuccess
	if isSuccess {
		data.Status = enum.ConfirmedStatus
		return
	}

	data.Status = enum.FailedStatus
}

func (e *ToncenterDataExtractor) formatRawAddress(workchainId int, rawAddress string) (string, error) {
	rawAddress = strings.Replace(rawAddress, fmt.Sprintf("%d:", workchainId), "", 1)
	bytes, err := hex.DecodeString(rawAddress)
	if err != nil {
		return "", fmt.Errorf("can not decode address: %s", err)
	}

	return address.NewAddress(17, byte(workchainId), bytes).Bounce(false).String(), nil
}
