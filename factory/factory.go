package factory

import (
	"fmt"
	"nexus-wallet/internal/env"
)

type ServiceFactory struct {
	env         *env.Env
	isDebugMode bool
}

func NewServiceFactory() (*ServiceFactory, error) {
	envData := &env.Env{}
	err := envData.Load()
	if err != nil {
		return nil, fmt.Errorf("env load failed: %s", err)
	}
	return &ServiceFactory{
		env:         envData,
		isDebugMode: envData.App.Mode == "DEV" || envData.App.Mode == "LOCAL",
	}, nil
}
