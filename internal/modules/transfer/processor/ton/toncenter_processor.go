package ton

import (
	"nexus-wallet/internal/app_error"
	processor_module "nexus-wallet/internal/modules/transfer/processor"
	"nexus-wallet/internal/shared/ton_message"
)

type processor struct {
	messageService *ton_message.TonMessageService
}

func NewProcessor(
	messageService *ton_message.TonMessageService,
) processor_module.TransferProcessor {
	return &processor{messageService: messageService}
}

func (p *processor) Transfer(message interface{}) (*processor_module.TransferOutput, *app_error.AppError) {
	hash, err := p.messageService.Send(message)
	if err != nil {
		return nil, err
	}

	return &processor_module.TransferOutput{Hash: hash}, nil
}
