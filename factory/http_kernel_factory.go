package factory

import (
	"fmt"
	"io"
	"nexus-wallet/api"
	"nexus-wallet/api/controllers"
	"nexus-wallet/api/controllers/address"
	"nexus-wallet/api/controllers/auth"
	"nexus-wallet/api/controllers/deposit"
	"nexus-wallet/api/controllers/estimate"
	"nexus-wallet/api/controllers/exchange"
	"nexus-wallet/api/controllers/mnemonic"
	"nexus-wallet/api/controllers/nft"
	profile_api "nexus-wallet/api/controllers/profile"
	"nexus-wallet/api/controllers/transaction"
	"nexus-wallet/api/controllers/transfer"
	"nexus-wallet/api/error_handler"
	api_logger "nexus-wallet/api/logger"
	"nexus-wallet/api/middlewares"
	"nexus-wallet/api/response"
	auth_module "nexus-wallet/internal/modules/auth"
	mnemonic_module "nexus-wallet/internal/modules/mnemonic"
)

func (f *ServiceFactory) CreateHttpKernel() (*api.Kernel, func() error, error) {
	connection, err := f.createConnection()
	if err != nil {
		return nil, nil, fmt.Errorf("create http kernel failed: %w", err)
	}
	sqlConnectionPool := connection.GetConnection()
	baseRepository, err := f.createBaseRepository(sqlConnectionPool)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create baseRepository: %s", err)
	}
	transactionManager := f.createTransactionManager(sqlConnectionPool)

	tokenPairService := f.createTokenPairService()

	coinService, err := f.createCoinService(baseRepository)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create coinService: %s", err)
	}
	authService, err := f.createAuthService(baseRepository, transactionManager, tokenPairService)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create authService: %s", err)
	}
	addressService, err := f.createAddressService(baseRepository, transactionManager)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create addressService: %s", err)
	}
	depositService, err := f.createDepositService(baseRepository)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create depositService: %s", err)
	}
	transferService, err := f.createTransferService(baseRepository)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create transferService: %s", err)
	}
	transactionService, err := f.createTransactionService(baseRepository)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create transactionService: %s", err)
	}
	mnemonicService := f.createMnemonicService(baseRepository, transactionManager)
	exchangeService, err := f.createExchangeService(baseRepository)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create exchangeService: %s", err)
	}
	nftService, err := f.createNftService(baseRepository, transactionManager)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create nftService: %s", err)
	}

	responseFactory, err := f.createHttpResponseFactory()
	if err != nil {
		return nil, nil, fmt.Errorf("can not create responseFactory: %s", err)
	}
	errorHandler, err := f.createHttpErrorHandler(responseFactory)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create errorHandler: %s", err)
	}
	accessTokenValidator := f.createHttpAccessTokenValidator(authService, errorHandler)
	mnemonicHashValidator := f.createMnemonicHashValidator(mnemonicService, errorHandler)
	rateLimiter := f.createHttpRateLimiter(f.env.App.RateLimit, responseFactory)

	authController := f.createAuthController(authService, errorHandler, responseFactory)
	telegramAuthController := f.createTelegramAuthController(authService, errorHandler, responseFactory)
	estimateController := f.createEstimateController(coinService, errorHandler, responseFactory)
	mnemonicController := f.createMnemonicController(mnemonicService, errorHandler, responseFactory)
	profileController := f.createProfileController(baseRepository, responseFactory)
	addressController := f.createAddressController(addressService, responseFactory, errorHandler)
	depositController := f.createDepositController(depositService, responseFactory, errorHandler)
	transferController := f.createTransferController(transferService, responseFactory, errorHandler)
	transactionController := f.createTransactionController(transactionService, responseFactory, errorHandler)
	exchangeController := f.createExchangeController(exchangeService, responseFactory, errorHandler)
	nftController := f.createNftController(nftService, responseFactory, errorHandler)

	apiLogger, err := f.createHttpApiLogger()
	if err != nil {
		return nil, nil, fmt.Errorf("can not create http engine logger, %s", err)
	}
	kernel := api.NewKernel(
		[]controllers.Router{
			address.NewRouter(addressController, accessTokenValidator, mnemonicHashValidator, rateLimiter),
			auth.NewRouter(authController, telegramAuthController, rateLimiter),
			estimate.NewRouter(estimateController, accessTokenValidator, rateLimiter),
			mnemonic.NewRouter(mnemonicController, accessTokenValidator, mnemonicHashValidator, rateLimiter),
			profile_api.NewRouter(profileController, accessTokenValidator, rateLimiter),
			deposit.NewRouter(depositController, accessTokenValidator, mnemonicHashValidator, rateLimiter),
			transfer.NewRouter(transferController, accessTokenValidator, mnemonicHashValidator, rateLimiter),
			transaction.NewRouter(transactionController, accessTokenValidator, mnemonicHashValidator, rateLimiter),
			exchange.NewRouter(exchangeController, accessTokenValidator, mnemonicHashValidator, rateLimiter),
			nft.NewRouter(nftController, accessTokenValidator, mnemonicHashValidator, rateLimiter),
		},
		apiLogger,
		errorHandler,
		api.Params{
			Port:            f.env.Http.Port,
			HeaderOriginURL: f.env.Http.HeaderOriginURL,
			DebugMode:       f.isDebugMode,
		},
	)

	onShutdown := func() error {
		return connection.GetConnection().Close()
	}
	return &kernel, onShutdown, nil
}

func (f *ServiceFactory) createHttpRateLimiter(
	maxRequestsPerSecond int,
	responseFactory *response.ResponseFactory,
) *middlewares.RateLimiter {
	return middlewares.NewRateLimiter(middlewares.Params{MaxRequests: float64(maxRequestsPerSecond)}, responseFactory)
}

func (f *ServiceFactory) createHttpAccessTokenValidator(
	service *auth_module.Service,
	errorHandler *error_handler.HttpErrorHandler,
) *middlewares.AccessTokenValidator {
	return middlewares.NewAccessTokenValidator(service, errorHandler)
}

func (f *ServiceFactory) createMnemonicHashValidator(
	service *mnemonic_module.Service,
	errorHandler *error_handler.HttpErrorHandler,
) *middlewares.MnemonicHashValidator {
	return middlewares.NewMnemonicHashValidator(service, errorHandler)
}

func (f *ServiceFactory) createHttpResponseFactory() (*response.ResponseFactory, error) {
	return response.NewResponseFactory(), nil
}

func (f *ServiceFactory) createHttpErrorHandler(
	responseFactory *response.ResponseFactory,
) (*error_handler.HttpErrorHandler, error) {
	newLogger, err := f.createLogger("http")
	if err != nil {
		return nil, fmt.Errorf("can not create http logger %s", err)
	}
	sentry, err := f.createSentry()
	if err != nil {
		return nil, fmt.Errorf("can not create create sentry: %s", err)
	}
	return error_handler.NewHttpErrorHandler(
		responseFactory,
		sentry,
		newLogger,
		error_handler.Params{IsDebugMode: f.isDebugMode},
	), nil
}

func (f *ServiceFactory) createHttpApiLogger() (io.Writer, error) {
	apiLogger, err := f.createLogger("http")
	if err != nil {
		return nil, fmt.Errorf("can not create http logger while creating apiLogger: %s", err)
	}
	return api_logger.NewApiLogger(apiLogger), nil
}
