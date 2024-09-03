package factory

import (
	"fmt"
	"github.com/pkg/errors"
	jwt "gitlab.com/golib4/jwt-resolver/jwt"
	"gitlab.com/golib4/telegram-query-serice/telegram_query"
	"nexus-wallet/internal/app_enum"
	address_module "nexus-wallet/internal/modules/address"
	"nexus-wallet/internal/modules/address/balance"
	balance_provider "nexus-wallet/internal/modules/address/balance/provider"
	ton_balance_provider "nexus-wallet/internal/modules/address/balance/provider/ton"
	trc20_balance_provider "nexus-wallet/internal/modules/address/balance/provider/trc20"
	"nexus-wallet/internal/modules/address/import"
	address_module_model "nexus-wallet/internal/modules/address/model/address"
	address_module_coin "nexus-wallet/internal/modules/address/model/coin"
	"nexus-wallet/internal/modules/address/model/mnemonic"
	auth_module "nexus-wallet/internal/modules/auth"
	"nexus-wallet/internal/modules/auth/factory/user"
	"nexus-wallet/internal/modules/auth/model/profile"
	auth_user "nexus-wallet/internal/modules/auth/model/user"
	auth_telegram "nexus-wallet/internal/modules/auth/telegram"
	"nexus-wallet/internal/modules/coin"
	coin_model "nexus-wallet/internal/modules/coin/model/coin"
	"nexus-wallet/internal/modules/coin/model/price"
	"nexus-wallet/internal/modules/coin/provider"
	"nexus-wallet/internal/modules/coin/provider/ton"
	"nexus-wallet/internal/modules/coin/provider/trc20"
	"nexus-wallet/internal/modules/deposit"
	"nexus-wallet/internal/modules/deposit/enum"
	"nexus-wallet/internal/modules/deposit/model/address"
	deposit_provider "nexus-wallet/internal/modules/deposit/provider"
	"nexus-wallet/internal/modules/deposit/provider/finch_pay"
	"nexus-wallet/internal/modules/deposit/provider/simple_swap"
	"nexus-wallet/internal/modules/exchange"
	exchange_address "nexus-wallet/internal/modules/exchange/model/address"
	exchange_coin_repository "nexus-wallet/internal/modules/exchange/model/coin"
	exchange_repository "nexus-wallet/internal/modules/exchange/model/exchange"
	"nexus-wallet/internal/modules/exchange/provider/change_hero"
	mnemonic_module "nexus-wallet/internal/modules/mnemonic"
	"nexus-wallet/internal/modules/mnemonic/model"
	"nexus-wallet/internal/modules/nft"
	nft_module_address "nexus-wallet/internal/modules/nft/model/address"
	nft_module_model "nexus-wallet/internal/modules/nft/model/nft"
	nft_syncer "nexus-wallet/internal/modules/nft/provider"
	nft_provider_ton "nexus-wallet/internal/modules/nft/provider/ton"
	trc20_nft_provider "nexus-wallet/internal/modules/nft/provider/trc20"
	"nexus-wallet/internal/modules/nft/sender"
	ton_nft_sender "nexus-wallet/internal/modules/nft/sender/ton"
	trc20_nft_sender "nexus-wallet/internal/modules/nft/sender/trc20"
	"nexus-wallet/internal/modules/notification"
	notification_module "nexus-wallet/internal/modules/notification/model/notification"
	notification_profile "nexus-wallet/internal/modules/notification/model/profile"
	notification_user "nexus-wallet/internal/modules/notification/model/user"
	notification_processor "nexus-wallet/internal/modules/notification/processor"
	"nexus-wallet/internal/modules/notification/processor/telegram"
	transaction_module "nexus-wallet/internal/modules/transaction"
	transaction_address "nexus-wallet/internal/modules/transaction/model/address"
	transaction_coin "nexus-wallet/internal/modules/transaction/model/coin"
	transaction2 "nexus-wallet/internal/modules/transaction/model/transaction"
	transaction_provider "nexus-wallet/internal/modules/transaction/provider"
	transaction_provider_ton "nexus-wallet/internal/modules/transaction/provider/ton"
	transaction_provider_trc20 "nexus-wallet/internal/modules/transaction/provider/trc20"
	"nexus-wallet/internal/modules/transfer"
	transfer_address "nexus-wallet/internal/modules/transfer/model/address"
	transfer_coin "nexus-wallet/internal/modules/transfer/model/coin"
	"nexus-wallet/internal/modules/transfer/processor"
	transfer_ton_processor "nexus-wallet/internal/modules/transfer/processor/ton"
	transfer_trc20_processor "nexus-wallet/internal/modules/transfer/processor/trc20"
	token_service "nexus-wallet/internal/shared/jwt"
	"nexus-wallet/internal/shared/ton_message"
	"nexus-wallet/pkg/repository"
	"nexus-wallet/pkg/transaction"
)

func (f *ServiceFactory) createTokenPairService() *token_service.TokenPairService {
	factory := jwt.NewJwtResolver(jwt.Params{
		SecretKey: []byte(f.env.App.SecretKey),
	})

	return token_service.NewTokenPairService(token_service.Params{
		AccessTokenLifetimeInMinutes:  f.env.Jwt.AccessTokenMaxLifeInMinutes,
		RefreshTokenLifetimeInMinutes: f.env.Jwt.RefreshTokenMaxLifeInMinutes,
	}, &factory)
}

func (f *ServiceFactory) createNftService(
	baseRepository *repository.BaseRepository,
	transactionManager transaction.Manager,
) (*nft.Service, error) {
	newLogger, err := f.createLogger("nft service")
	if err != nil {
		return nil, fmt.Errorf("can not create logger while creating nftService: %s", err)
	}

	addressRepository := nft_module_address.NewRepository(baseRepository)
	nftRepository := nft_module_model.NewRepository(baseRepository)

	tonconsoleClient, err := f.createTonconsoleClient()
	if err != nil {
		return nil, fmt.Errorf("can not create client for tonconsole in nftService: %s", err)
	}

	tonMessageService, err := f.createTonMessageService()
	if err != nil {
		return nil, fmt.Errorf("can not create tonMessageService in nftService: %s", err)
	}

	nftMessageBuilders := make(map[app_enum.Network]sender.NftMessageBuilder)
	nftMessageBuilders[app_enum.TonNetwork] = ton_nft_sender.NewBuilder(tonMessageService)
	nftMessageBuilders[app_enum.Trc20Network] = trc20_nft_sender.NewBuilder()

	nftSenders := make(map[app_enum.Network]sender.NftSender)
	nftSenders[app_enum.TonNetwork] = ton_nft_sender.NewSender(tonMessageService)
	nftSenders[app_enum.Trc20Network] = trc20_nft_sender.NewSender()

	nftProviders := make(map[app_enum.Network]nft_syncer.NftProvider)
	nftProviders[app_enum.TonNetwork] = nft_provider_ton.NewProvider(tonconsoleClient)
	nftProviders[app_enum.Trc20Network] = trc20_nft_provider.NewProvider()

	syncer := nft_syncer.NewSyncer(
		f.createSyncManager(f.env.Sync.NftIntervalMinutes),
		newLogger,
		addressRepository,
		nftRepository,
		nftProviders,
		transactionManager,
		nft_syncer.Params{MaxAttempts: 3},
	)

	return nft.NewService(nftMessageBuilders, nftSenders, addressRepository, nftRepository, syncer), nil
}

func (f *ServiceFactory) createExchangeService(
	baseRepository *repository.BaseRepository,
) (*exchange.Service, error) {
	priceRepository := exchange_coin_repository.NewPriceRepository(baseRepository)

	changeHeroClient, err := f.createChangeHeroClient()
	if err != nil {
		return nil, errors.Errorf("create change hero client error: %s", err.Error())
	}

	return exchange.NewService(
		change_hero.NewChangeHeroProvider(changeHeroClient, priceRepository),
		exchange_address.NewRepository(baseRepository),
		exchange_repository.NewRepository(baseRepository),
		exchange_coin_repository.NewRepository(baseRepository),
	), nil
}

func (f *ServiceFactory) createTransactionService(
	baseRepository *repository.BaseRepository,
) (*transaction_module.Service, error) {
	newLogger, err := f.createLogger("transaction service")
	if err != nil {
		return nil, fmt.Errorf("can not create logger while creating transactionService: %s", err)
	}

	coinRepository := transaction_coin.NewRepository(baseRepository)
	addressRepository := transaction_address.NewRepository(baseRepository)
	transactionRepository := transaction2.NewTransactionRepository(baseRepository)

	tronGridDataExtractor, err := transaction_provider_trc20.NewTronDataExtractor(coinRepository)
	if err != nil {
		return nil, fmt.Errorf("can not create trongrid data extractor: %s", err)
	}
	tronGridClient, err := f.createTrongridClient()
	if err != nil {
		return nil, fmt.Errorf("can not create client for trongrid in transactionService: %s", err)
	}
	tronscanapiClient, err := f.createTronscanapiClient()
	if err != nil {
		return nil, fmt.Errorf("can not create client for tronscanapi in transactionService: %s", err)
	}
	newTronGridLogger, err := f.createLogger("tronGrid transactions provider")
	if err != nil {
		return nil, fmt.Errorf("can not create trongridLogger while creating transactionService: %s", err)
	}

	tonCenterExtractor := transaction_provider_ton.NewToncenterDataExtractor()
	tonCenterClient, err := f.createToncenterClient()
	if err != nil {
		return nil, fmt.Errorf("can not create client for toncenter in transactionService: %s", err)
	}
	newTonCenterLogger, err := f.createLogger("tonCenter transactions provider")
	if err != nil {
		return nil, fmt.Errorf("can not create toncenterLogger while creating transactionService: %s", err)
	}

	providers := make(map[app_enum.Network]transaction_provider.Provider)
	providers[app_enum.Trc20Network] = transaction_provider_trc20.NewTronProvider(tronGridClient, tronGridDataExtractor, tronscanapiClient, newTronGridLogger)
	providers[app_enum.TonNetwork] = transaction_provider_ton.NewToncenterProvider(tonCenterClient, tonCenterExtractor, newTonCenterLogger)

	return transaction_module.NewService(
		transactionRepository,
		coinRepository,
		transaction_provider.NewSyncer(
			transactionRepository,
			addressRepository,
			coinRepository,
			newLogger,
			providers,
			f.createSyncManager(f.env.Sync.TransactionsIntervalMinutes),
		),
	), nil
}

func (f *ServiceFactory) createMnemonicService(
	baseRepository *repository.BaseRepository,
	transactionManager transaction.Manager,
) *mnemonic_module.Service {
	mnemonicRepository := model.NewRepository(baseRepository)
	return mnemonic_module.NewService(
		mnemonicRepository,
		f.createNotficationCreator(baseRepository, transactionManager),
	)
}

func (f *ServiceFactory) createTransferService(
	baseRepository *repository.BaseRepository,
) (*transfer.Service, error) {
	coinRepository := transfer_coin.NewRepository(baseRepository)
	addressRepository := transfer_address.NewRepository(baseRepository)
	addressCoinRepository := transfer_address.NewAddressCoinRepository(baseRepository)

	trongridGrpcClient, err := f.createTrongridGrpcClient()
	if err != nil {
		return nil, fmt.Errorf("can not create client for tron grpc in transferService: %s", err)
	}
	tonMessageService, err := f.createTonMessageService()
	if err != nil {
		return nil, fmt.Errorf("can not create tonMessageService in transferService: %s", err)
	}

	processors := make(map[app_enum.Network]processor.TransferProcessor)

	processors[app_enum.TonNetwork] = transfer_ton_processor.NewProcessor(tonMessageService)
	processors[app_enum.Trc20Network] = transfer_trc20_processor.NewProcessor(trongridGrpcClient)

	builders := make(map[app_enum.Network]processor.TransferMessageBuilder)
	builders[app_enum.TonNetwork] = transfer_ton_processor.NewBuilder(tonMessageService)
	builders[app_enum.Trc20Network] = transfer_trc20_processor.NewBuilder(trongridGrpcClient)

	newLogger, err := f.createLogger("transfer service")
	if err != nil {
		return nil, fmt.Errorf("can not create http logger while creating transfer service: %s", err)
	}

	return transfer.NewService(processors, builders, coinRepository, addressRepository, addressCoinRepository, f.createCacher(), newLogger), nil
}

func (f *ServiceFactory) createDepositService(
	baseRepository *repository.BaseRepository,
) (*deposit.Service, error) {
	simpleSwapClient, err := f.createSimpleSwapClient()
	if err != nil {
		return nil, fmt.Errorf("can not create simpleSwapClient: %s", err)
	}
	finchPayClient, err := f.createFinchTechClient()
	if err != nil {
		return nil, fmt.Errorf("can not create finchPayClient: %s", err)
	}

	providers := make(map[enum.ProviderName]deposit_provider.RedirectProvider)
	providers[enum.FinchPayProviderName] = finch_pay.NewFinchPayProvider(
		finch_pay.Params{
			WidgetHost: f.env.Integrations.FinchPay.WidgetHost,
			PartnerId:  f.env.Integrations.FinchPay.PartnerId,
			SecretKey:  f.env.Integrations.FinchPay.SecretKey,
		},
		finchPayClient,
	)
	providers[enum.SimpleSwapProviderName] = simple_swap.NewSimpleSwapProvider(simpleSwapClient)

	return deposit.NewService(providers, address.NewRepository(baseRepository)), nil
}

func (f *ServiceFactory) createAuthService(
	baseRepository *repository.BaseRepository,
	transactionManager transaction.Manager,
	tokenPairService *token_service.TokenPairService,
) (*auth_module.Service, error) {
	userRepository := auth_user.NewRepository(baseRepository)
	profileRepository := profile.NewRepository(baseRepository)
	telegramProfileRepository := profile.NewTelegramProfileRepository(baseRepository)

	newLogger, err := f.createLogger("telegram auth service")
	if err != nil {
		return nil, fmt.Errorf("can not create http logger while creating authService: %s", err)
	}
	telegramClient, err := f.createTelegramClient()
	if err != nil {
		return nil, fmt.Errorf("can not create telegram client: %s", err)
	}
	telegramUserRepository := auth_user.NewTelegramUserRepository(userRepository)
	authTelegramService := auth_telegram.NewService(
		telegram_query.NewQueryService(telegram_query.Params{
			BotToken:                       f.env.Telegram.BotToken,
			TelegramQueryLifetimeInSeconds: f.env.Telegram.QueryLifetimeInSeconds,
		}),
		user.NewTelegramUserFactory(
			telegramClient,
			profileRepository,
			telegramProfileRepository,
			telegramUserRepository,
			transactionManager,
			newLogger,
			f.env.Telegram.BotToken,
		),
		telegramUserRepository,
	)
	return auth_module.NewService(authTelegramService, tokenPairService), nil
}

func (f *ServiceFactory) createCoinService(baseRepository *repository.BaseRepository) (*coin.Service, error) {
	tonhubapiClient, err := f.createTonhubClient()
	if err != nil {
		return nil, fmt.Errorf("can not create tonhubapiClient for coin service: %s", err)
	}
	coingeckoClient, err := f.createCoingeckoClient()
	if err != nil {
		return nil, fmt.Errorf("can not create coingeckoClient for coin service: %s", err)
	}

	providers := []provider.Provider{
		ton.NewTonProvider(tonhubapiClient),
		trc20.NewTetherProvider(coingeckoClient),
		trc20.NewTronProvider(coingeckoClient),
	}
	newLogger, err := f.createLogger("coin service")
	if err != nil {
		return nil, fmt.Errorf("can not create http logger while creating nftController: %s", err)
	}

	priceRepository := price.NewRepository(baseRepository)
	coinRepository := coin_model.NewRepository(baseRepository)
	return coin.NewService(
		priceRepository,
		coinRepository,
		provider.NewSyncer(
			priceRepository,
			coinRepository,
			newLogger,
			providers,
			provider.Params{FrequencyInHours: f.env.Sync.PricesFrequencyInHours},
		),
	), nil
}

func (f *ServiceFactory) createAddressService(
	baseRepository *repository.BaseRepository,
	transactionManager transaction.Manager,
) (*address_module.Service, error) {
	newLogger, err := f.createLogger("address balance service")
	if err != nil {
		return nil, fmt.Errorf("can not create http logger while creating addressService: %s", err)
	}

	toncenterClient, err := f.createToncenterClient()
	if err != nil {
		return nil, fmt.Errorf("can not create client for toncenter in addressService: %s", err)
	}
	trongridGrpcClient, err := f.createTrongridGrpcClient()
	if err != nil {
		return nil, fmt.Errorf("can not create client for tron grpc in addressService: %s", err)
	}
	tronscanapiClient, err := f.createTronscanapiClient()
	if err != nil {
		return nil, fmt.Errorf("can not create client for tronscanapi in addressService: %s", err)
	}

	providers := make(map[app_enum.Network][]balance_provider.Provider)
	providers[app_enum.TonNetwork] = []balance_provider.Provider{
		ton_balance_provider.NewToncenterProvider(toncenterClient),
	}
	providers[app_enum.Trc20Network] = []balance_provider.Provider{
		trc20_balance_provider.NewTrongridGrpcProvider(trongridGrpcClient),
		trc20_balance_provider.NewTronscanapiProvider(tronscanapiClient),
	}

	addressCoinRepository := address_module_model.NewAddressCoinRepository(baseRepository)
	coinRepository := address_module_coin.NewCoinRepository(baseRepository)
	addressRepository := address_module_model.NewAddressRepository(baseRepository)
	mnemonicRepository := mnemonic.NewRepository(baseRepository)

	return address_module.NewService(
		addressCoinRepository,
		addressRepository,
		_import.NewImporter(addressRepository, addressCoinRepository, mnemonicRepository, coinRepository, transactionManager),
		coinRepository,
		address_module_coin.NewPriceRepository(baseRepository),
		mnemonicRepository,
		balance.NewSyncer(
			addressCoinRepository,
			addressRepository,
			coinRepository,
			newLogger,
			providers,
			balance.Params{MaxAttempts: 3},
			f.createSyncManager(f.env.Sync.BalanceIntervalMinutes),
		),
	), nil
}

func (f *ServiceFactory) createNotificationService(
	baseRepository *repository.BaseRepository,
) (*notification.Service, error) {
	telegramClient, err := f.createTelegramClient()
	if err != nil {
		return nil, fmt.Errorf("can not create telegram client: %s", err)
	}

	processors := make(map[app_enum.ProfileType]notification_processor.NotificationProcessor)
	processors[app_enum.TelegramProfileType] = telegram.NewTelegramProcessor(telegramClient, notification_user.NewTelegramUserRepository(baseRepository))

	return notification.NewService(
		processors,
		notification_profile.NewRepository(baseRepository),
		notification_module.NewGlobalNotificationRepository(baseRepository),
		notification_module.NewGlobalNotificationTranslationRepository(baseRepository),
		notification_module.NewGlobalProcessedNotificationRepository(baseRepository),
		notification_module.NewPersonalNotificationRepository(baseRepository),
		notification_module.NewPersonalNotificationTranslationRepository(baseRepository),
		notification.Params{IsFake: f.env.App.Mode == "LOCAL"},
	), nil
}

func (f *ServiceFactory) createTonMessageService() (*ton_message.TonMessageService, error) {
	toncenterClient, err := f.createToncenterClient()
	if err != nil {
		return nil, fmt.Errorf("can not create client for toncenter in transferService: %s", err)
	}

	return ton_message.NewTonMessageService(toncenterClient), nil
}
