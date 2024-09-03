package factory

import (
	"nexus-wallet/api/controllers/address"
	"nexus-wallet/api/controllers/auth"
	"nexus-wallet/api/controllers/deposit"
	"nexus-wallet/api/controllers/estimate"
	"nexus-wallet/api/controllers/exchange"
	"nexus-wallet/api/controllers/mnemonic"
	"nexus-wallet/api/controllers/nft"
	"nexus-wallet/api/controllers/profile"
	"nexus-wallet/api/controllers/transaction"
	"nexus-wallet/api/controllers/transfer"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/api/response"
	address_module "nexus-wallet/internal/modules/address"
	auth_module "nexus-wallet/internal/modules/auth"
	"nexus-wallet/internal/modules/coin"
	deposit_module "nexus-wallet/internal/modules/deposit"
	exchange_module "nexus-wallet/internal/modules/exchange"
	mnemonic_module "nexus-wallet/internal/modules/mnemonic"
	nft_module "nexus-wallet/internal/modules/nft"
	profile_module "nexus-wallet/internal/modules/profile"
	profile_model "nexus-wallet/internal/modules/profile/model"
	transaction_module "nexus-wallet/internal/modules/transaction"
	transfer_module "nexus-wallet/internal/modules/transfer"
	"nexus-wallet/pkg/repository"
)

func (f *ServiceFactory) createTransactionController(
	service *transaction_module.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *transaction.TransactionController {
	return transaction.NewTransactionController(service, responseFactory, errorHandler)
}

func (f *ServiceFactory) createExchangeController(
	service *exchange_module.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *exchange.ExchangeController {
	return exchange.NewExchangeController(service, responseFactory, errorHandler)
}

func (f *ServiceFactory) createNftController(
	service *nft_module.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *nft.NftController {
	return nft.NewNftController(responseFactory, errorHandler, service)
}

func (f *ServiceFactory) createTransferController(
	service *transfer_module.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *transfer.TransferController {
	return transfer.NewTransferController(service, responseFactory, errorHandler)
}

func (f *ServiceFactory) createDepositController(
	service *deposit_module.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *deposit.DepositController {
	return deposit.NewDepositController(service, responseFactory, errorHandler)
}

func (f *ServiceFactory) createAddressController(
	service *address_module.Service,
	responseFactory *response.ResponseFactory,
	errorHandler *error_handler.HttpErrorHandler,
) *address.AddressController {
	return address.NewAddressController(service, responseFactory, errorHandler)
}

func (f *ServiceFactory) createTelegramAuthController(
	authService *auth_module.Service,
	errorHandler *error_handler.HttpErrorHandler,
	responseFactory *response.ResponseFactory,
) *auth.TelegramAuthController {
	return auth.NewTelegramAuthController(authService, errorHandler, responseFactory)
}

func (f *ServiceFactory) createAuthController(
	authService *auth_module.Service,
	errorHandler *error_handler.HttpErrorHandler,
	responseFactory *response.ResponseFactory,
) *auth.AuthController {
	return auth.NewAuthController(authService, errorHandler, responseFactory)
}

func (f *ServiceFactory) createEstimateController(
	coinService *coin.Service,
	errorHandler *error_handler.HttpErrorHandler,
	responseFactory *response.ResponseFactory,
) *estimate.EstimateController {
	return estimate.NewEstimateController(coinService, responseFactory, errorHandler)
}

func (f *ServiceFactory) createProfileController(
	baseRepository *repository.BaseRepository,
	responseFactory *response.ResponseFactory,
) *profile.ProfileController {
	service := profile_module.NewService(profile_model.NewRepository(baseRepository))
	return profile.NewProfileController(service, responseFactory)
}

func (f *ServiceFactory) createMnemonicController(
	service *mnemonic_module.Service,
	errorHandler *error_handler.HttpErrorHandler,
	responseFactory *response.ResponseFactory,
) *mnemonic.MnemonicController {
	return mnemonic.NewMnemonicController(service, responseFactory, errorHandler)
}
