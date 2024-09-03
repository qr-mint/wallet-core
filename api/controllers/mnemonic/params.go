package mnemonic

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/internal/modules/mnemonic"
)

type ValidateHashRequest struct {
	Value string `json:"value"`
}

func (ValidateHashRequest) createInputFromRequest(context *gin.Context) (*mnemonic.GetMnemonicIdInput, *app_error.AppError) {
	body := ValidateHashRequest{}
	if err := context.BindJSON(&body); err != nil {
		return nil, app_error.InvalidDataError(errors.New("Invalid request"))
	}
	if body.Value == "" {
		return nil, app_error.InvalidDataError(errors.New("`value` cannot be empty"))
	}

	return &mnemonic.GetMnemonicIdInput{Hash: body.Value}, nil
}

type ValidateHashResponse struct {
	IsValid bool `json:"is_valid"`
}

type GenerateRequest struct {
}

func (GenerateRequest) createInputFromRequest(context *gin.Context) mnemonic.GenerateInput {
	return mnemonic.GenerateInput{
		UserId: context.GetInt64("userId"),
	}
}

type GenerateResponse struct {
	Value string `json:"value"`
}

type GetNamesRequest struct {
}

func (GetNamesRequest) createInputFromRequest(context *gin.Context) mnemonic.GetNamesInput {
	return mnemonic.GetNamesInput{
		UserId: context.GetInt64("userId"),
	}
}

type GetNamesResponseItem struct {
	Value string `json:"value"`
}

type GetNamesResponse struct {
	Names []GetNamesResponseItem `json:"items"`
}

func (GetNamesResponse) fillFromOutput(output mnemonic.GetNamesOutput) GetNamesResponse {
	var namesData []GetNamesResponseItem
	for _, name := range output.Items {
		namesData = append(namesData, GetNamesResponseItem{Value: name.Name})
	}

	return GetNamesResponse{Names: namesData}
}

type UpdateNameRequest struct {
	Value string `json:"value"`
}

func (UpdateNameRequest) createInputFromRequest(context *gin.Context) (*mnemonic.UpdateNameInput, *app_error.AppError) {
	body := UpdateNameRequest{}
	if err := context.BindJSON(&body); err != nil {
		return nil, app_error.InvalidDataError(errors.New("Invalid request"))
	}
	if len(body.Value) > 20 {
		return nil, app_error.InvalidDataError(errors.New("`value` length too long"))
	}

	if body.Value == "" {
		return nil, app_error.InvalidDataError(errors.New("param `value` is required"))
	}
	return &mnemonic.UpdateNameInput{Name: body.Value, MnemonicId: context.GetInt64("mnemonicId")}, nil
}
