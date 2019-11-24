package logging

import (
	"errors"
	"testing"

	"go.uber.org/zap"
)

func TestInterfaces(t *testing.T) {
	InitLogger()
	Debug("msg")
	Info("msg")
	Warn("msg")
	Error("中文")
	Error("with error field", zap.Error(errors.New("iamerror")))
}
