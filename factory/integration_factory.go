package factory

import (
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/golib4/coingecko-client/coingecko"
	"gitlab.com/golib4/toncenter-client/toncenter"
	"gitlab.com/golib4/tonconsole-client/tonconsole"
	"gitlab.com/golib4/tonhubapi-client/tonhubapi"
	"gitlab.com/golib4/trongrid-client/trongrid"
	"gitlab.com/golib4/tronscanapi-client/tronscanapi"
	"google.golang.org/grpc"
	"nexus-wallet/pkg/change_hero"
	finch_pay_client "nexus-wallet/pkg/finch_pay"
	"nexus-wallet/pkg/simple_swap"
)

func (f *ServiceFactory) createTelegramClient() (*tgbotAPI.BotAPI, error) {
	telegramClient, err := tgbotAPI.NewBotAPI(f.env.Telegram.BotToken)
	if err != nil {
		return nil, fmt.Errorf("can not create telegram client: %s", err)
	}

	return telegramClient, nil
}

func (f *ServiceFactory) createSimpleSwapClient() (*simple_swap.Client, error) {
	httpClient, err := f.createHttpClient(f.env.Integrations.SimpleSwap.Host, 10)
	if err != nil {
		return nil, fmt.Errorf("can not create http client while creating simpleSwapClient: %s", err)
	}

	return simple_swap.NewClient(httpClient, simple_swap.Params{ApiKey: f.env.Integrations.SimpleSwap.ApiKey}), nil
}

func (f *ServiceFactory) createFinchTechClient() (*finch_pay_client.Client, error) {
	httpClient, err := f.createHttpClient(f.env.Integrations.FinchPay.Host, 10)
	if err != nil {
		return nil, fmt.Errorf("can not create http client while creating simpleSwapClient: %s", err)
	}

	return finch_pay_client.NewClient(httpClient, finch_pay_client.Params{ApiKey: f.env.Integrations.FinchPay.PartnerId}), nil
}

func (f *ServiceFactory) createChangeHeroClient() (*change_hero.Client, error) {
	httpClient, err := f.createHttpClient(f.env.Integrations.ChangeHero.Host, 10)
	if err != nil {
		return nil, fmt.Errorf("can not create http client while creating changeHeroClient: %s", err)
	}

	return change_hero.NewClient(httpClient, change_hero.Params{ApiKey: f.env.Integrations.ChangeHero.Apikey}), nil
}

func (f *ServiceFactory) createTonconsoleClient() (*tonconsole.Client, error) {
	httpClient, err := f.createHttpClient(f.env.Integrations.Tonconsole.Host, 10)
	if err != nil {
		return nil, fmt.Errorf("can not create http client while creating tonconsoleClient: %s", err)
	}
	newClient := tonconsole.NewClient(httpClient)

	return &newClient, nil
}

func (f *ServiceFactory) createToncenterClient() (*toncenter.Client, error) {
	httpClient, err := f.createHttpClient(f.env.Integrations.Toncenter.Host, 10)
	if err != nil {
		return nil, fmt.Errorf("can not create http client while creating toncenterClient: %s", err)
	}
	newClient := toncenter.NewClient(httpClient, f.env.Integrations.Toncenter.Key)

	return &newClient, nil
}

func (f *ServiceFactory) createTonhubClient() (*tonhubapi.Client, error) {
	httpClient, err := f.createHttpClient(f.env.Integrations.Tonhub.ConnectHost, 10)
	if err != nil {
		return nil, fmt.Errorf("can not create http client while creating tonhubClient: %s", err)
	}
	newClient := tonhubapi.NewClient(httpClient)

	return &newClient, nil
}

func (f *ServiceFactory) createCoingeckoClient() (*coingecko.Client, error) {
	httpClient, err := f.createHttpClient(f.env.Integrations.Coingecko.Host, 10)
	if err != nil {
		return nil, fmt.Errorf("can not create http client for coingecko: %s", err)
	}

	newClient := coingecko.NewClient(httpClient)
	return &newClient, nil
}

func (f *ServiceFactory) createTrongridGrpcClient() (*client.GrpcClient, error) {
	trongridGrpcClient := client.NewGrpcClient("")
	err := trongridGrpcClient.SetAPIKey(f.env.Integrations.Trongrid.Key)
	if err != nil {
		return nil, fmt.Errorf("can not set api tron grpc client: %s", err)
	}
	err = trongridGrpcClient.Start(grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("can not start tron grpc client: %s", err)
	}

	return trongridGrpcClient, nil
}

func (f *ServiceFactory) createTrongridClient() (*trongrid.Client, error) {
	httpClient, err := f.createHttpClient(f.env.Integrations.Trongrid.Host, 10)
	if err != nil {
		return nil, fmt.Errorf("can not create http client while creating trongridClient: %s", err)
	}
	newClient := trongrid.NewClient(httpClient, f.env.Integrations.Trongrid.Key)

	return &newClient, nil
}

func (f *ServiceFactory) createTronscanapiClient() (*tronscanapi.Client, error) {
	httpClient, err := f.createHttpClient(f.env.Integrations.Tronscanapi.Host, 10)
	if err != nil {
		return nil, fmt.Errorf("can not create http client while creating tronscanapiClient: %s", err)
	}
	newClient := tronscanapi.NewClient(httpClient, f.env.Integrations.Tronscanapi.Key)

	return &newClient, nil
}
