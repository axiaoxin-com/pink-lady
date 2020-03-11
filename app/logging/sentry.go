package logging

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// SentryClient godoc
var SentryClient *sentry.Client

// InitSentry init sentry
func InitSentry() error {
	options := sentry.ClientOptions{
		Dsn:              viper.GetString("server.sentrydsn"),
		Debug:            viper.GetBool("server.mode"),
		AttachStacktrace: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					// You have access to the original Request
					Debug("Sentry BeforeSend req:", req)
				}
			}
			Debug("Sentry BeforeSend event:", event)
			return event
		},
	}
	err := sentry.Init(options)

	SentryClient, err = sentry.NewClient(options)

	return errors.Wrap(err, "init sentry error")
}
