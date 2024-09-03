package nft

import (
	"fmt"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/app_util"
	"nexus-wallet/internal/modules/nft/model/address"
	"nexus-wallet/internal/modules/nft/model/nft"
	"nexus-wallet/internal/modules/nft/provider"
	"nexus-wallet/internal/modules/nft/sender"
)

type Service struct {
	nftMessageBuilders map[app_enum.Network]sender.NftMessageBuilder
	nftSenders         map[app_enum.Network]sender.NftSender
	addressRepository  *address.Repository
	nftRepository      *nft.Repository
	syncer             *provider.Syncer
}

func NewService(
	nftMessageBuilders map[app_enum.Network]sender.NftMessageBuilder,
	nftSenders map[app_enum.Network]sender.NftSender,
	addressRepository *address.Repository,
	nftRepository *nft.Repository,
	syncer *provider.Syncer,
) *Service {
	return &Service{
		nftMessageBuilders: nftMessageBuilders,
		nftSenders:         nftSenders,
		addressRepository:  addressRepository,
		nftRepository:      nftRepository,
		syncer:             syncer,
	}
}

func (s Service) List(input ListInput) (*ListOutput, *app_error.AppError) {
	err := s.syncer.Sync(input.MnemonicId)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not sync nft: %s", err))
	}

	nftList, err := s.nftRepository.FindMany(nft.FindManyOptions{MnemonicId: input.MnemonicId, Limit: input.Limit, Offset: input.Offset}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get nft list: %s", err))
	}
	addressList, err := s.addressRepository.FindMany(address.FindManyOptions{MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get nft list: %s", err))
	}

	result, err := ListOutput{}.fillFromModel(nftList, addressList)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not fill from model nft list: %s", err))
	}

	return result, nil
}

func (s Service) Get(input GetInput) (*GetOutput, *app_error.AppError) {
	nftModel, err := s.nftRepository.Find(nft.FindOptions{MnemonicId: input.MnemonicId, Id: input.NftId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not get nft: %s", err))
	}
	if nftModel == nil {
		return nil, app_error.ResourceNotFoundError(errors.Errorf("nft %d not found", input.NftId))
	}
	addressList, err := s.addressRepository.FindMany(address.FindManyOptions{MnemonicId: input.MnemonicId}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get nft list: %s", err))
	}

	result, err := GetOutput{}.fillFromModel(*nftModel, addressList[nftModel.AddressId])
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not fill from model nft get: %s", err))
	}

	return result, nil
}

func (s Service) BuildSendMessage(input BuildSendMessageInput) (*BuildSendMessageOutput, *app_error.AppError) {
	isValidAddress, err := app_util.IsValidAddress(input.ToAddress, input.Network)
	if !isValidAddress {
		return nil, app_error.InvalidDataError(fmt.Errorf("incorrect to address: %s", input.ToAddress))
	}
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not validate address: %s", err))
	}

	fromAddress, err := s.addressRepository.Find(address.FindOptions{MnemonicId: input.MnemonicId, Network: input.Network}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not find address: %s", err))
	}
	if fromAddress == nil {
		return nil, app_error.InternalError(errors.Errorf("address of network %s, not found for mnemonic id: %d", input.Network, input.MnemonicId))
	}
	if fromAddress.Address == input.ToAddress {
		return nil, app_error.InvalidDataError(errors.New("addressTo can not be equal to addressFrom"))
	}
	nftData, err := s.nftRepository.FindOne(input.NftId, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not find nft: %s", err))
	}
	if nftData == nil {
		return nil, app_error.InvalidDataError(errors.New("nft is not found by id"))
	}

	nftMessageBuilder, nftMessageBuilderFounded := s.nftMessageBuilders[input.Network]
	if !nftMessageBuilderFounded {
		return nil, app_error.InternalError(errors.Errorf("not found nft message builder by network: %s", input.Network))
	}

	toTransaction, buildNftSendErr := nftMessageBuilder.BuildSendNftMessage(sender.BuildNftMessageInput{
		NftAddress:  nftData.Address,
		ToAddress:   input.ToAddress,
		FromAddress: fromAddress.Address,
		Options:     sender.BuildMessageWalletInputOptions{Version: input.Version, PublicKey: input.PublicKey},
	})
	if buildNftSendErr != nil {
		return nil, buildNftSendErr
	}

	return &BuildSendMessageOutput{Message: toTransaction.Message}, nil
}

func (s Service) Send(input SendInput) (*SendOutput, *app_error.AppError) {
	nftSender, nftSenderFounded := s.nftSenders[input.Network]
	if !nftSenderFounded {
		return nil, app_error.InternalError(errors.Errorf("not found nft sender by network: %s", input.Network))
	}

	result, sendErr := nftSender.Send(input.Message)
	if sendErr != nil {
		return nil, sendErr
	}

	nftModel, err := s.nftRepository.Find(nft.FindOptions{MnemonicId: input.MnemonicId, Id: input.NftId}, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not get nft model while send: %s", err))
	}
	if nftModel == nil {
		return nil, app_error.InvalidDataError(errors.Errorf("nft model is not found by id: %d", input.NftId))
	}
	err = s.nftRepository.Delete(nftModel, nil)
	if err != nil {
		return nil, app_error.InternalError(errors.Errorf("can not delete nft model while send: %s", err))
	}

	return &SendOutput{Hash: result.Hash}, nil
}
