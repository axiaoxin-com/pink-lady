package utils

import (
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func InitSentry() error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   viper.GetString("server.sentrydsn"),
		Debug: viper.GetBool("server.mode"),
	})
	return errors.Wrap(err, "init sentry error")
}
