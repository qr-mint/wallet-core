package cron

import (
	"fmt"
	"github.com/gitploy-io/cronexpr"
	"gitlab.com/golib4/logger/logger"
	"nexus-wallet/cron/error_handler"
	"nexus-wallet/cron/jobs"
	"nexus-wallet/internal/app_error"
	"runtime/debug"
	"sync"
	"time"
)

type Kernel struct {
	runners      []jobs.Runner
	errorHandler *error_handler.CronErrorHandler
	logger       logger.Logger
}

func NewKernel(
	runners []jobs.Runner,
	errorHandler *error_handler.CronErrorHandler,
	logger logger.Logger,
) Kernel {
	return Kernel{
		runners:      runners,
		errorHandler: errorHandler,
		logger:       logger,
	}
}

func (k Kernel) Run() {
	for _, runner := range k.runners {
		go func() {
			defer k.registerRecover(runner)

			var mu sync.Mutex
			for {
				nextTime := cronexpr.MustParse(runner.GetPattern()).Next(time.Now())
				time.Sleep(time.Until(nextTime))

				mu.Lock()
				k.logger.Info("starting cron runner ...")
				runner.Run()
				mu.Unlock()
			}
		}()

		fmt.Printf("\n started cron runner %T \n", runner)
	}
}

func (k Kernel) registerRecover(runner jobs.Runner) {
	if err := recover(); err != nil {
		k.errorHandler.Handle(
			app_error.InternalError(fmt.Errorf("panic: %s. stack: %s", err, debug.Stack())),
			fmt.Sprintf("%T", runner),
		)
	}
}
