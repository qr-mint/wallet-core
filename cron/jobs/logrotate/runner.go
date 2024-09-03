package logrotate

import (
	"gitlab.com/golib4/logger/logger"
)

type Runner struct {
	logger logger.Logger
}

func NewRunner(
	logger logger.Logger,
) *Runner {
	return &Runner{
		logger: logger,
	}
}

func (r Runner) Run() {
	count := 0
	for count < 1000 {
		r.logger.Clear(1000)

		count++
	}
}

func (Runner) GetPattern() string {
	return "0 0 * * 0" // every week
}
