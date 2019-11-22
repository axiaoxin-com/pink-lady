package utils

import (
	"net/http"

	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func InitSentry() error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              viper.GetString("server.sentrydsn"),
		Debug:            viper.GetBool("server.mode"),
		AttachStacktrace: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			logger := logging.Logger.Sugar()
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					// You have access to the original Request
					logger.Debug("Sentry BeforeSend req:", req)
				}
			}
			logger.Debug("Sentry BeforeSend event:", event)
			return event
		},
	})
	return errors.Wrap(err, "init sentry error")
}
