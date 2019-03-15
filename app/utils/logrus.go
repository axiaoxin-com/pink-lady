package utils

import (
	"io"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Logger logrus entry
var Logger *logrus.Entry

// InitLogger set logrus logger options
// logLevel set level which will be logged, values: debug, info(default), warning, error, fatal, panic
// logFormatter set log format, values: text(default), json
func InitLogger(output io.Writer, logLevel string, logFormatter string) {
	level, err := logrus.ParseLevel(logLevel)
	if err == nil {
		logrus.SetLevel(level)
	}

	logger := logrus.New()
	if strings.ToLower(logFormatter) == "json" {
		logger.Formatter = &logrus.JSONFormatter{}
	} else {
		logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	}

	logger.Out = output
	Logger = logrus.NewEntry(logger)
	Logger = Logger.WithFields(logrus.Fields{
		"PID": syscall.Getpid(),
	})
}
