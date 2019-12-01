package logging

import (
	"syscall"
	"testing"

	"github.com/spf13/viper"
)

func TestInitLogger(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Fatal(err)
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
		false,
		false,
	)
	if logger == nil || err != nil {
		t.Fatal(err)
	}
}
