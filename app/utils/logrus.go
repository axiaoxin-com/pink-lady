package utils

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// init logrus
// logLevel set level which will be logged, values: debug, info(default), warning, error, fatal, panic
// logFormatter set log format, values: text(default), json
func InitLogrus(logLevel string, logFormatter string) {
	level, err := logrus.ParseLevel(logLevel)
	if err == nil {
		logrus.SetLevel(level)
	}

	if strings.ToLower(logFormatter) == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	logrus.SetOutput(os.Stdout)

}
