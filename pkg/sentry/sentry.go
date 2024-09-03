package sentry

import (
	"fmt"
	"github.com/getsentry/sentry-go"
)

type SentryParams struct {
	DebugMode bool
	SentryDsn string
}

type Sentry struct {
	params SentryParams
}

func NewSentry(params SentryParams) (*Sentry, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   params.SentryDsn,
		Debug: params.DebugMode,
	})
	if err != nil {
		return nil, fmt.Errorf("sentry initialization failed: %s", err)
	}

	return &Sentry{params: params}, nil
}

func (Sentry) HandleError(err error) {
	go sentry.CaptureException(err)
}
