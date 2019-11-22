package utils

import (
	"testing"

	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

func TestInitSentry(t *testing.T) {
	logging.InitLogger()
	err := InitSentry()
	if err != nil {
		t.Error(err)
	}
	sentry.CaptureException(errors.New("pink-lady-testing"))
}
