package factory

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"gitlab.com/golib4/http-client/http"
	"gitlab.com/golib4/logger/logger"
	"nexus-wallet/internal/shared/notification_creator"
	"nexus-wallet/internal/shared/notification_creator/model"
	"nexus-wallet/internal/shared/sync"
	"nexus-wallet/pkg/cache"
	"nexus-wallet/pkg/repository"
	"nexus-wallet/pkg/sentry"
	"nexus-wallet/pkg/transaction"
)

func (f *ServiceFactory) createLogger(source string) (logger.Logger, error) {
	if f.env.App.Mode != "PROD" && f.env.App.Mode != "DEV" {
		return logger.NewOutputLogger(), nil
	}

	newClient, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%s", f.env.Elastic.Host, f.env.Elastic.Port)),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", f.env.Elastic.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("can not create elastic client: %s", err)
	}
	return logger.NewElasticLogger(newClient, logger.Params{
		Source:      source,
		Environment: f.env.App.Mode,
		IndexName:   "logs",
	}), nil
}

func (f *ServiceFactory) createHttpClient(host string, timeoutInSeconds int) (*http.Client, error) {
	apiLogger, err := f.createLogger("api")
	if err != nil {
		return nil, err
	}
	newClient := http.NewClient(host, apiLogger, http.Params{LogInfo: true, TimeoutInSeconds: timeoutInSeconds})
	return &newClient, err
}

func (f *ServiceFactory) createSentry() (*sentry.Sentry, error) {
	return sentry.NewSentry(sentry.SentryParams{
		DebugMode: f.isDebugMode,
		SentryDsn: f.env.Sentry.Dsn,
	})
}

func (f *ServiceFactory) createCacher() cache.Cacher {
	return cache.NewRedisCacher(cache.Params{
		Host:     f.env.Redis.Host,
		Port:     f.env.Redis.Port,
		Password: f.env.Redis.Password,
		Database: f.env.Redis.Database,
	})
}

func (f *ServiceFactory) createNotficationCreator(
	baseRepository *repository.BaseRepository,
	transactionManager transaction.Manager,
) *notification_creator.Creator {
	return notification_creator.NewCreator(
		model.NewRepository(baseRepository),
		transactionManager,
	)
}

func (f *ServiceFactory) createSyncManager(syncIntervalInMinutes int) *sync.Manager {
	return sync.NewManager(
		f.createCacher(),
		sync.Params{SyncIntervalInMinutes: syncIntervalInMinutes},
	)
}
