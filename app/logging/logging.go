// Package logging provides ...
package logging

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Debug 记录debug级别的日志
func Debug(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Debug(msg, fields...)
}

// Info 记录info级别的日志
func Info(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Info(msg, fields...)
}

// Warn 记录warn级别的日志
func Warn(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Warn(msg, fields...)
}

// Error 记录Error级别的日志，如果fields中有zap.Error则上报的sentry
func Error(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Error(msg, fields...)
	for _, field := range fields {
		if field.Type == zapcore.ErrorType {
			sentry.CaptureException(field.Interface.(error))
		}
	}
}
