package notify_personal

import (
	"fmt"
	"nexus-wallet/cron/error_handler"
	"nexus-wallet/internal/modules/notification"
)

type Runner struct {
	service      *notification.Service
	errorHandler *error_handler.CronErrorHandler
}

func NewRunner(
	service *notification.Service,
	errorHandler *error_handler.CronErrorHandler,
) *Runner {
	return &Runner{
		service:      service,
		errorHandler: errorHandler,
	}
}

func (r Runner) Run() {
	err := r.service.PersonalNotify()
	if err != nil {
		r.errorHandler.Handle(err, fmt.Sprintf("%T", r))
	}
}

func (Runner) GetPattern() string {
	return "* * * * *" // every 1 minute
}
