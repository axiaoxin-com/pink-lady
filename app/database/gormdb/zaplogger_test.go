package gormdb

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(zap.L())
	if logger.zap == nil {
		t.Fatal("logger.zap is nil")
	}
}

func TestPrint(t *testing.T) {

	logger := NewLogger(zap.NewExample())
	logger.Print()
	logger.Print("-", "1", time.Duration(2), "3", "4", int64(5))
	logger.Print("sql", "1", time.Duration(2), "3", "4", int64(5))
}
