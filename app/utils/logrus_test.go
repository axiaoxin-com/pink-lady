package utils

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestInitLogrus(t *testing.T) {
	InitLogrus(os.Stdout, "debug", "text")
	if logrus.GetLevel().String() != "debug" {
		t.Error("set level failure")
	}
	InitLogrus(os.Stdout, "error", "text")
	if logrus.GetLevel().String() != "error" {
		t.Error("set level failure")
	}
}
