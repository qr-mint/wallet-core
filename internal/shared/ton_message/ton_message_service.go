package ton_message

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/pkg/errors"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"gitlab.com/golib4/toncenter-client/toncenter"
	"nexus-wallet/internal/app_error"
	"time"
)

type TonMessageService struct {
	client *toncenter.Client
}

func NewTonMessageService(client *toncenter.Client) *TonMessageService {
	return &TonMessageService{client: client}
}

type SignedMessage struct {
	DestinationAddress string
	StateInit          *tlb.StateInit
	Body               *cell.Cell
	Signature          string
}

func (s TonMessageService) Send(message interface{}) (string, *app_error.AppError) {
	signedMessage, isValidType := message.(SignedMessage)
	if !isValidType {
		return "", app_error.InternalError(errors.Errorf("message must be of type SignerMessage, %T provided", message))
	}

	if signedMessage.Signature == "" || signedMessage.DestinationAddress == "" || signedMessage.Body == nil {
		return "", app_error.InvalidDataError(errors.New("invalid message provided"))
	}

	signatureHex, err := common.Hex2Bytes(signedMessage.Signature)
	if err != nil {
		return "", app_error.InvalidDataError(errors.Errorf("can not convert signature to hex: %s", err))
	}

	externalMessage, err := tlb.ToCell(&tlb.ExternalMessage{
		DstAddr:   address.MustParseAddr(signedMessage.DestinationAddress),
		StateInit: signedMessage.StateInit,
		Body:      cell.BeginCell().MustStoreSlice(signatureHex, 512).MustStoreBuilder(signedMessage.Body.ToBuilder()).EndCell(),
	})
	if err != nil {
		return "", app_error.InternalError(errors.Errorf("can not convert message to cell: %s", err))
	}

	response, err := s.client.SendMessage(
		toncenter.SendMessageRequest{
			Message: base64.StdEncoding.EncodeToString(externalMessage.ToBOCWithFlags(false)),
		},
	)
	if err != nil {
		return "", app_error.InternalError(errors.Errorf("can not send message to toncenter: %s", err))
	}

	return response.Result.Hash, nil
}

type BuiltMessage struct {
	DestinationAddress string
	StateInit          *tlb.StateInit
	Body               *cell.Cell
}

func (s TonMessageService) BuildExternalMessage(
	fromAddress string,
	versionInt uint8,
	publicKeyString string,
	internalMessage wallet.Message,
) (*BuiltMessage, *app_error.AppError) {
	if !s.isWalletVersionSupported(versionInt) {
		return nil, app_error.InvalidDataError(fmt.Errorf("ton send is not yet supported for version %d", versionInt))
	}

	addressInfo, err := s.client.GetAddressInfo(toncenter.GetAddressInfoRequest{Address: fromAddress})
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("failed to get ton account seqno from toncenter: %s", err))
	}

	body, err := s.buildExternalMessageBody(internalMessage, uint64(addressInfo.Seqno))
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("build ton message err: %s", err))
	}

	if addressInfo.Status != "active" {
		stateInit, err := s.getStateInit(publicKeyString)
		if err != nil {
			return nil, err
		}

		return &BuiltMessage{DestinationAddress: fromAddress, StateInit: stateInit, Body: body}, nil

	}

	return &BuiltMessage{DestinationAddress: fromAddress, Body: body}, nil
}

func (TonMessageService) buildExternalMessageBody(message wallet.Message, seqNo uint64) (*cell.Cell, error) {
	cellInternalMessage, err := tlb.ToCell(message.InternalMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to convert internal message to cell: %s", err)
	}

	messageTTL := 10 * time.Minute

	return cell.BeginCell().MustStoreUInt(uint64(wallet.DefaultSubwallet), 32).
		MustStoreUInt(uint64(time.Now().Add(messageTTL).UTC().Unix()), 32).
		MustStoreUInt(seqNo, 32).
		MustStoreInt(0, 8).
		MustStoreUInt(uint64(message.Mode), 8).MustStoreRef(cellInternalMessage).
		EndCell(), nil
}

func (TonMessageService) getStateInit(publicKeyString string) (*tlb.StateInit, *app_error.AppError) {
	if publicKeyString == "" {
		return nil, app_error.InvalidDataError(errors.New("publicKey is required for ton network"))
	}
	pubBytes, err := hex.DecodeString(publicKeyString)
	if err != nil {
		return nil, app_error.InvalidDataError(errors.Errorf("failed to decode ton public key: %s", err))
	}
	publicKey := ed25519.PublicKey(pubBytes)

	stateInit, err := wallet.GetStateInit(publicKey, wallet.V4R2, wallet.DefaultSubwallet)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to get state init: %s", err))
	}

	return stateInit, nil
}

func (TonMessageService) isWalletVersionSupported(versionInt uint8) bool {
	version := wallet.Version(versionInt)
	if version != wallet.V3R2 && version != wallet.V3R1 && version != wallet.V4R2 && version != wallet.V4R1 {
		return false
	}

	return true
}
