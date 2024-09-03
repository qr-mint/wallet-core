package error_handler

import (
	"gitlab.com/golib4/logger/logger"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/pkg/sentry"
)

type CronErrorHandler struct {
	logger logger.Logger
	sentry *sentry.Sentry
}

func NewCronErrorHandler(logger logger.Logger, sentry *sentry.Sentry) *CronErrorHandler {
	return &CronErrorHandler{
		logger: logger,
		sentry: sentry,
	}
}

func (h *CronErrorHandler) Handle(err *app_error.AppError, cronJobName string) {
	switch err.Code {
	default:
		h.sentry.HandleError(err.Error)
		h.logger.Errorf("error in cronjob `%s`: %s", cronJobName, err.Error.Error())
	}
}
