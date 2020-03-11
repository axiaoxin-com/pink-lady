// Package logging provides ...
package logging

// Debug 记录debug级别的日志
// logging.Debug("abc", 123)
func Debug(args ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Debug(args...)
}

// Info 记录info级别的日志
func Info(args ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Info(args...)
}

// Warn 记录warn级别的日志
func Warn(args ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Warn(args...)
}

// Error 记录Error级别的日志
func Error(args ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Error(args...)
}

// Debugf 模板字符串记录debug级别的日志
// logging.Debugf("str:%s", "abd")
func Debugf(template string, args ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Debugf(template, args...)
}

// Infof 模板字符串记录info级别的日志
func Infof(template string, args ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Infof(template, args...)
}

// Warnf 模板字符串记录warn级别的日志
func Warnf(template string, args ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Warnf(template, args...)
}

// Errorf 模板字符串记录debug级别的日志
func Errorf(template string, args ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Errorf(template, args...)
}

// Debugw kv记录debug级别的日志
// logging.Debugw("msg", "k1", "v1", "k2", "v2")
func Debugw(msg string, keysAndValues ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Debugw(msg, keysAndValues...)
}

// Infow kv记录info级别的日志
func Infow(msg string, keysAndValues ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Infow(msg, keysAndValues...)
}

// Warnw kv记录warn级别的日志
func Warnw(msg string, keysAndValues ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Warnw(msg, keysAndValues...)
}

// Errorw kv记录error级别的日志
func Errorw(msg string, keysAndValues ...interface{}) {
	slogger := Logger.Sugar()
	defer slogger.Sync()
	slogger.Errorw(msg, keysAndValues...)
}
