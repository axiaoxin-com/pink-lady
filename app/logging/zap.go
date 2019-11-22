// Package logger provides ...
package logging

import (
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger global zap logger with pid field
	Logger *zap.Logger
	// DefaultZapOutPaths zap日志默认输出位置
	DefaultZapOutPaths = []string{"stderr"}
)

// InitLogger init the global zap logger
func InitLogger() error {
	var err error
	Logger, err = NewLogger(
		viper.GetString("logger.level"),
		viper.GetString("logger.format"),
		viper.GetStringSlice("logger.outputPaths"),
		map[string]interface{}{
			"pid": syscall.Getpid(),
		},
	)
	return err
}

// NewLogger return a new zap logger
func NewLogger(level, format string, outputPaths []string, initialFields map[string]interface{}) (*zap.Logger, error) {
	cfg := zap.Config{}
	// 设置level 默认为info
	switch strings.ToLower(level) {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "dpanic":
		cfg.Level = zap.NewAtomicLevelAt(zap.DPanicLevel)
	case "panic":
		cfg.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		cfg.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	// 设置encoding 默认为json
	switch strings.ToLower(format) {
	case "console":
		cfg.Encoding = "console"
	default:
		cfg.Encoding = "json"
	}
	// 设置output 默认全部输出到stderr
	switch len(outputPaths) {
	case 0:
		cfg.OutputPaths = DefaultZapOutPaths
		cfg.ErrorOutputPaths = DefaultZapOutPaths
	default:
		cfg.OutputPaths = outputPaths
		cfg.ErrorOutputPaths = outputPaths
	}
	// 设置disablecaller
	cfg.DisableCaller = false
	// 设置disablestacktrace
	cfg.DisableStacktrace = false
	// 设置InitialFields
	cfg.InitialFields = initialFields

	// 设置encoderConfig
	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Sampling实现了日志的流控功能，或者叫采样配置，主要有两个配置参数，Initial和Thereafter，实现的效果是在1s的时间单位内，如果某个日志级别下同样内容的日志输出数量超过了Initial的数量，那么超过之后，每隔Thereafter的数量，才会再输出一次。是一个对日志输出的保护功能。
	cfg.Sampling = &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

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

// Error 记录Error级别的日志
func Error(msg string, fields ...zap.Field) {
	defer Logger.Sync()
	Logger.Error(msg, fields...)
}

// Clone return the global Logger copy
func CloneLogger() *zap.Logger {
	copy := *Logger
	return &copy
}
