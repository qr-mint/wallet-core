package trc20

import (
	"encoding/hex"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/shengdoushi/base58"
	trc20decoder "gitlab.com/golib4/trc20decoder"
	"gitlab.com/golib4/trongrid-client/trongrid"
	"gitlab.com/golib4/tronscanapi-client/tronscanapi"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/modules/transaction/enum"
	"nexus-wallet/internal/modules/transaction/model/coin"
	"nexus-wallet/internal/modules/transaction/provider"
	"strconv"
	"time"

	"strings"
)

type TronDataExtractor struct {
	coins          []*coin.Coin
	coinRepository *coin.Repository
}

func NewTronDataExtractor(coinRepository *coin.Repository) (*TronDataExtractor, error) {
	return &TronDataExtractor{
		coins:          nil,
		coinRepository: coinRepository,
	}, nil
}

func (e *TronDataExtractor) Extract(transaction trongrid.Transaction, account string) (*provider.BlockchainTransactionData, error) {
	if e.coins == nil {
		coins, err := e.coinRepository.FindMany(coin.FindManyOptions{Network: app_enum.Trc20Network}, nil)
		if err != nil {
			return nil, fmt.Errorf("can not find trc20 coin data: %s", err)
		}

		e.coins = coins
	}

	if e.needToSkip(transaction) {
		return nil, nil
	}

	data := &provider.BlockchainTransactionData{}

	e.setMainDataFromTransaction(transaction, data)
	e.setStatus(transaction, data)
	e.setType(account, data)
	if err := e.setAmountAndToAddress(transaction, data); err != nil {
		return nil, fmt.Errorf("can not set amount and to address: %s", err)
	}
	if err := e.setCoinName(transaction, data); err != nil {
		return nil, fmt.Errorf("can not set coin name to data: %s", err)
	}

	return data, nil
}

func (e *TronDataExtractor) ExtractToken(transfer tronscanapi.TokenTransfer, account string) (*provider.BlockchainTransactionData, error) {
	status := enum.FailedStatus
	if transfer.FinalResult == "SUCCESS" {
		status = enum.ConfirmedStatus
	}
	transactionType := enum.InType
	if transfer.FromAddress == account {
		transactionType = enum.OutType
	}
	intAmount, err := strconv.Atoi(transfer.Quant)
	if err != nil {
		return nil, fmt.Errorf("failed to convert tron token amount: %s", err)
	}
	if transfer.TokenInfo.TokenName != "Tether USD" {
		return nil, nil
	}

	return &provider.BlockchainTransactionData{
		From:      transfer.FromAddress,
		To:        transfer.ToAddress,
		Amount:    int64(intAmount),
		Hash:      transfer.TransactionID,
		Type:      transactionType,
		CoinName:  app_enum.TetherCoinName,
		Network:   app_enum.Trc20Network,
		Status:    status,
		CreatedAt: time.Unix(0, transfer.BlockTs*int64(time.Millisecond)),
	}, nil

}

func (e *TronDataExtractor) needToSkip(transaction trongrid.Transaction) bool {
	if len(transaction.RawData.Contract) > 1 {
		return true
	}

	notTransfer := transaction.RawData.Contract[0].Parameter.Value.Data != nil &&
		!trc20decoder.IsTransfer(*transaction.RawData.Contract[0].Parameter.Value.Data)

	return notTransfer
}

func (e *TronDataExtractor) setMainDataFromTransaction(transaction trongrid.Transaction, data *provider.BlockchainTransactionData) {
	data.From = address.HexToAddress(transaction.RawData.Contract[0].Parameter.Value.OwnerAddress).String()
	data.Hash = transaction.TxID
	data.CreatedAt = time.Unix(0, transaction.RawData.Timestamp*int64(time.Millisecond))
	data.Network = app_enum.Trc20Network
}

func (e *TronDataExtractor) setType(account string, data *provider.BlockchainTransactionData) {
	transactionType := enum.InType
	if data.From == account {
		transactionType = enum.OutType
	}

	data.Type = transactionType
}

func (e *TronDataExtractor) setAmountAndToAddress(transaction trongrid.Transaction, data *provider.BlockchainTransactionData) error {
	if transaction.RawData.Contract[0].Parameter.Value.Amount != nil && transaction.RawData.Contract[0].Parameter.Value.ToAddress != nil {
		data.Amount = int64(*transaction.RawData.Contract[0].Parameter.Value.Amount)
		data.To = address.HexToAddress("0x" + *transaction.RawData.Contract[0].Parameter.Value.ToAddress).String()

		return nil
	}

	if transaction.RawData.Contract[0].Parameter.Value.Data != nil {
		decodedData, err := trc20decoder.Decode(*transaction.RawData.Contract[0].Parameter.Value.Data)
		if err != nil {
			return fmt.Errorf("can not decode trc20 data: %s", err)
		}

		data.Amount = decodedData.Amount
		data.To = address.HexToAddress("0x" + decodedData.ToAddress).String()

		return nil
	}

	return fmt.Errorf("can not extract data from transaction %s with type %s", transaction.TxID, transaction.RawData.Contract[0].Type)
}

func (e *TronDataExtractor) setStatus(transaction trongrid.Transaction, data *provider.BlockchainTransactionData) {
	for _, ret := range transaction.Ret {
		if ret.ContractRet == "SUCCESS" {
			data.Status = enum.ConfirmedStatus
			return
		}
		if ret.ContractRet == "FAILED" {
			data.Status = enum.FailedStatus
			return
		}
	}

	data.Status = enum.FailedStatus
}

func (e *TronDataExtractor) setCoinName(transaction trongrid.Transaction, data *provider.BlockchainTransactionData) error {
	for _, c := range e.coins {
		if c.Address == nil {
			continue
		}

		hexCoinAddress, err := e.transformAddressToHex(*c.Address)
		if err != nil {
			return fmt.Errorf("error transforming address to hex: %s", err)
		}

		if !strings.Contains(transaction.RawDataHex, hexCoinAddress) {
			continue
		}

		data.CoinName = c.Name

		return nil
	}

	data.CoinName = app_enum.TronCoinName

	return nil
}

func (e *TronDataExtractor) transformAddressToHex(address string) (string, error) {
	bytes, err := base58.Decode(address, base58.BitcoinAlphabet)
	if err != nil {
		return "", fmt.Errorf("error wile decoding address: %s", err)
	}

	return hex.EncodeToString(bytes[:len(bytes)-4]), nil
}
