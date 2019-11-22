package logging

import (
	"errors"
	"syscall"
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Error(err)
	}
}

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(
		viper.GetString("logger.level"),
		viper.GetString("logger.format"),
		viper.GetStringSlice("logger.outputPaths"),
		map[string]interface{}{
			"pid": syscall.Getpid(),
		},
	)
	if logger == nil || err != nil {
		t.Error(err)
	}

}

func TestInterfaces(t *testing.T) {
	InitLogger()
	Debug("msg")
	Info("msg")
	Warn("msg")
	Error("中文")
	Error("with error field", zap.Error(errors.New("iamerror")))
}
