package ton

import (
	"nexus-wallet/internal/app_error"
	sender_module "nexus-wallet/internal/modules/nft/sender"
	"nexus-wallet/internal/shared/ton_message"
)

type sender struct {
	messageService *ton_message.TonMessageService
}

func NewSender(messageService *ton_message.TonMessageService) sender_module.NftSender {
	return &sender{messageService: messageService}
}

func (p *sender) Send(message interface{}) (*sender_module.SendOutput, *app_error.AppError) {
	hash, err := p.messageService.Send(message)
	if err != nil {
		return nil, err
	}

	return &sender_module.SendOutput{Hash: hash}, nil
}
