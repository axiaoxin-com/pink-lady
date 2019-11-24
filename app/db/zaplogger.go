package db

import (
	"time"

	"go.uber.org/zap"
)

// Logger 使用zap作为gorm的logger
type Logger struct {
	zap *zap.Logger
}

// NewLogger 创建gorm的zaplogger
func NewLogger(logger *zap.Logger) Logger {
	return Logger{zap: logger}
}

// Print 实现gorm logger接口
func (l Logger) Print(values ...interface{}) {
	if len(values) < 2 {
		return
	}

	switch values[0] {
	case "sql":
		l.zap.Debug("gorm.debug.sql",
			zap.String("query", values[3].(string)),
			zap.Any("values", values[4]),
			zap.Duration("duration", values[2].(time.Duration)),
			zap.Int64("affected-rows", values[5].(int64)),
			zap.String("source", values[1].(string)), // if AddCallerSkip(6) is well defined, we can safely remove this field
		)
	default:
		l.zap.Debug("gorm.debug.other",
			zap.Any("values", values[2:]),
			zap.String("source", values[1].(string)), // if AddCallerSkip(6) is well defined, we can safely remove this field
		)
	}
}
