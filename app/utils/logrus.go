package utils

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// init logrus
// logLevel set level which will be logged, values: debug, info(default), warning, error, fatal, panic
// logFormatter set log format, values: text(default), json
// out set where the log will be output, values: stderr, stdout(default)
func InitLogrus(logLevel string, logFormatter string, out string) {
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
	if strings.ToLower(out) == "stderr" {
		logrus.SetOutput(os.Stderr)
	}

}
