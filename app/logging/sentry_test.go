package logging

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

func TestInitSentry(t *testing.T) {
	InitLogger()
	err := InitSentry()
	if err != nil {
		t.Fatal(err)
	}
	sentry.CaptureException(errors.New("pink-lady-testing"))
}
