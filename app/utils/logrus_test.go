package utils

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestInitLogrus(t *testing.T) {
	InitLogrus("debug", "text")
	if logrus.GetLevel().String() != "debug" {
		t.Error("set level failure")
	}
	InitLogrus("error", "text")
	if logrus.GetLevel().String() != "error" {
		t.Error("set level failure")
	}
}
