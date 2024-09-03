package mnemonic

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/golib4/mnemonic-generator/mnemonic"
	"nexus-wallet/internal/app_enum"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/mnemonic/model"
	"nexus-wallet/internal/shared/notification_creator"
)

type Service struct {
	repository          *model.Repository
	notificationCreator *notification_creator.Creator
}

func NewService(
	repository *model.Repository,
	notificationCreator *notification_creator.Creator,
) *Service {
	return &Service{
		repository:          repository,
		notificationCreator: notificationCreator,
	}
}

func (s Service) GetMnemonicId(input GetMnemonicIdInput) (*GetMnemonicIdOutput, *app_error.AppError) {
	mnemonicData, err := s.repository.Find(model.FindOptions{Hash: input.Hash}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not get mnemonic hash: %s", err))
	}
	if mnemonicData == nil {
		return &GetMnemonicIdOutput{Id: 0}, nil
	}

	return &GetMnemonicIdOutput{Id: mnemonicData.Id}, nil
}

func (s Service) Generate(input GenerateInput) (*GenerateOutput, *app_error.AppError) {
	mnemonicString, err := mnemonic.GenerateMnemonic("", mnemonic.Size12)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not generate mnemonic: %s", err))
	}

	texts := make(map[app_enum.Language]string)
	texts[app_enum.EngLanguage] = "‼️ Не забудьте сохранить вашу секретную фразу."
	texts[app_enum.RuLanguage] = "‼️ Don't forget to save your secret phrase."

	appErr := s.notificationCreator.Create(notification_creator.CreateInput{
		UserId: input.UserId,
		Texts:  texts,
	})
	if appErr != nil {
		return nil, appErr
	}

	return &GenerateOutput{Mnemonic: mnemonicString}, nil
}

func (s Service) UpdateName(input UpdateNameInput) *app_error.AppError {
	mnemonicData, err := s.repository.FindOne(input.MnemonicId, nil)
	if err != nil {
		return app_error.InternalError(errors.Errorf("failed to find mnemonic in wallet service: %s", err))
	}
	if mnemonicData == nil {
		return app_error.InvalidDataError(errors.Errorf("invald mnemonic id provided: %d, does not exists", input.MnemonicId))
	}

	mnemonicData.Name = input.Name
	err = s.repository.Save(mnemonicData, nil)
	if err != nil {
		return app_error.InternalError(errors.Errorf("can not save mnemonic in wallet service: %s", err))
	}

	return nil
}

func (s Service) GetNames(data GetNamesInput) (*GetNamesOutput, *app_error.AppError) {
	items, err := s.repository.FindMany(model.FindManyOptions{UserId: data.UserId}, nil)
	if err != nil {
		return nil, app_error.InternalError(fmt.Errorf("can not find many wallet names in wallet service: %s", err))
	}

	var outputItems []GetNamesOutputItem
	for _, item := range items {
		outputItems = append(outputItems, GetNamesOutputItem{Name: item.Name})
	}

	return &GetNamesOutput{outputItems}, nil
}
