package enum

import (
	"nexus-wallet/internal/app_enum/utils"
)

type ProviderName string

const (
	SimpleSwapProviderName ProviderName = "simple_swap"
	FinchPayProviderName   ProviderName = "finch_pay"
)

func ToProviderName(value string) *ProviderName {
	if !utils.AssertInArray(value, []string{string(SimpleSwapProviderName), string(FinchPayProviderName)}) {
		return nil
	}

	providerName := ProviderName(value)
	return &providerName
}
