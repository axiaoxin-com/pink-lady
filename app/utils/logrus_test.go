package utils

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestInitLogger(t *testing.T) {
	InitLogger(os.Stdout, "debug", "text")
	if logrus.GetLevel().String() != "debug" {
		t.Error("set level failure")
	}
	if Logger == nil {
		t.Error("Logger init failed")
	}
	InitLogger(os.Stdout, "error", "json")
	if logrus.GetLevel().String() != "error" {
		t.Error("set level failure")
	}
}
