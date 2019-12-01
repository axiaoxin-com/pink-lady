package utils

import (
	"testing"

	"pink-lady/app/logging"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

func TestInitSentry(t *testing.T) {
	logging.InitLogger()
	err := InitSentry()
	if err != nil {
		t.Fatal(err)
	}
	sentry.CaptureException(errors.New("pink-lady-testing"))
}
