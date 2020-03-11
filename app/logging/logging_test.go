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

	Debugf("%s", "asf")
	Infof("%s", "asf")
	Warnf("%s", "asf")
	Errorf("%s", "asf")

	Debugw("asf", "k1", "v1")
	Infow("asf", "k1", "v1")
	Warnw("asf", "k1", "v1")
	Errorw("asf", "k1", "v1")
}
