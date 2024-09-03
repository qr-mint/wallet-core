package ton

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"gitlab.com/golib4/coins/coins"
	"nexus-wallet/internal/app_error"
	processor_module "nexus-wallet/internal/modules/transfer/processor"
	"nexus-wallet/internal/shared/ton_message"
)

type builder struct {
	messageService *ton_message.TonMessageService
}

func NewBuilder(messageService *ton_message.TonMessageService) processor_module.TransferMessageBuilder {
	return &builder{
		messageService: messageService,
	}
}

func (b builder) BuildMessage(input processor_module.BuildMessageInput) (*processor_module.BuildMessageOutput, *app_error.AppError) {
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

	return &processor_module.BuildMessageOutput{Message: externalMessage}, nil
}

func (b builder) BuildTokenMessage(input processor_module.BuildTokenMessageInput) (*processor_module.BuildMessageOutput, *app_error.AppError) {
	return nil, app_error.IllegalOperationError(errors.New("ton BuildTokenMessage is not implemented"))
}

func (builder) buildInternalMessage(input processor_module.BuildMessageInput) (*wallet.Message, error) {
	var body *cell.Cell
	var err error
	if input.Options.Comment != "" {
		body, err = wallet.CreateCommentCell(input.Options.Comment)
	}
	if err != nil {
		return nil, errors.Errorf("failed to create ton comment cell: %s", err)
	}

	return &wallet.Message{
		Mode: 1 + 2,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      false,
			DstAddr:     address.MustParseAddr(input.ToAddress),
			Amount:      tlb.FromNanoTON(coins.MustFromDecimal(fmt.Sprintf("%f", input.Amount), 9).Nano()),
			Body:        body,
		},
	}, nil
}
