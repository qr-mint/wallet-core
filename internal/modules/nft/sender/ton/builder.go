package ton

import (
	"nexus-wallet/internal/app_error"
	sender_module "nexus-wallet/internal/modules/nft/sender"
	"nexus-wallet/internal/shared/ton_message"

	"github.com/pkg/errors"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type builder struct {
	messageService *ton_message.TonMessageService
}

func NewBuilder(messageService *ton_message.TonMessageService) sender_module.NftMessageBuilder {
	return builder{messageService: messageService}
}

func (b builder) BuildSendNftMessage(input sender_module.BuildNftMessageInput) (*sender_module.BuildMessageOutput, *app_error.AppError) {
	internalMessage, err := b.buildInternalMessage(input)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("failed to build ton internal message: %s", err))
	}

	externalMessage, buildErr := b.messageService.BuildExternalMessage(
		input.FromAddress,
		input.Options.Version,
		input.Options.PublicKey,
		*internalMessage,
	)
	if buildErr != nil {
		return nil, buildErr
	}

	return &sender_module.BuildMessageOutput{Message: externalMessage}, nil

}

func (builder) buildInternalMessage(input sender_module.BuildNftMessageInput) (*wallet.Message, error) {
	payloadForward := cell.BeginCell().EndCell()

	newOwner := address.MustParseAddr(input.ToAddress)
	respTo := address.MustParseAddr(input.FromAddress)

	body, err := tlb.ToCell(nft.TransferPayload{
		QueryID:             0,
		NewOwner:            newOwner,
		ResponseDestination: respTo,
		CustomPayload:       nil,
		ForwardAmount:       tlb.MustFromTON("0.02"),
		ForwardPayload:      payloadForward,
	})
	if err != nil {
		return nil, errors.Errorf("can not build ton nft transfer payload: %s", err)
	}

	return &wallet.Message{
		Mode: 1 + 2,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      false,
			DstAddr:     address.MustParseAddr(input.NftAddress),
			Amount:      tlb.MustFromTON("0.065"),
			Body:        body,
		},
	}, nil
}
