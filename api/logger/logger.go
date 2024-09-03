package logger

import (
	"gitlab.com/golib4/logger/logger"
	"io"
)

type apiLogger struct {
	logger logger.Logger
}

func NewApiLogger(logger logger.Logger) io.Writer {
	return apiLogger{logger: logger}
}

func (l apiLogger) Write(p []byte) (n int, err error) {
	l.logger.Info(string(p))

	return 0, nil
}
