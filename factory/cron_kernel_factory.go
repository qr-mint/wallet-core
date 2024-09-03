package factory

import (
	"fmt"
	"nexus-wallet/cron"
	"nexus-wallet/cron/error_handler"
	"nexus-wallet/cron/jobs"
	"nexus-wallet/cron/jobs/logrotate"
	"nexus-wallet/cron/jobs/notify_global"
	"nexus-wallet/cron/jobs/notify_personal"
	"nexus-wallet/cron/jobs/sync_coin_prices"
)

func (f *ServiceFactory) CreateCronKernel() (*cron.Kernel, func() error, error) {
	connection, err := f.createConnection()
	if err != nil {
		return nil, nil, fmt.Errorf("create http kernel failed: %w", err)
	}
	sqlConnectionPool := connection.GetConnection()
	baseRepository, err := f.createBaseRepository(sqlConnectionPool)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create baseRepository: %s", err)
	}

	notificationService, err := f.createNotificationService(baseRepository)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create notificationService: %s", err)
	}
	coinService, err := f.createCoinService(baseRepository)
	if err != nil {
		return nil, nil, fmt.Errorf("can not create coinService: %s", err)
	}

	errorHandler, err := f.createCronErrorHandler()
	if err != nil {
		return nil, nil, fmt.Errorf("can not create cron errorHandler: %s", err)
	}

	logger, err := f.createLogger("cron")
	if err != nil {
		return nil, nil, fmt.Errorf("can not create logger for cron: %s", err)
	}

	kernel := cron.NewKernel(
		[]jobs.Runner{
			notify_global.NewRunner(notificationService, errorHandler),
			notify_personal.NewRunner(notificationService, errorHandler),
			sync_coin_prices.NewRunner(coinService),
			logrotate.NewRunner(logger),
		},
		errorHandler,
		logger,
	)
	onShutdown := func() error {
		return connection.GetConnection().Close()
	}
	return &kernel, onShutdown, nil
}

func (f *ServiceFactory) createCronErrorHandler() (*error_handler.CronErrorHandler, error) {
	logger, err := f.createLogger("cron")
	if err != nil {
		return nil, fmt.Errorf("can not create cron logger: %s", err)
	}
	sentry, err := f.createSentry()
	if err != nil {
		return nil, fmt.Errorf("can not create create sentry: %s", err)
	}
	return error_handler.NewCronErrorHandler(logger, sentry), nil
}
