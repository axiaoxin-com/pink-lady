package utils

import (
	"net/http"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/sirupsen/logrus"
)

func EndlessServe(bind string, handler http.Handler) {
	server := endless.NewServer(bind, handler)
	server.BeforeBegin = func(bind string) {
		logrus.Infof("Server is listening and serving HTTP on %s (pids: %d)", bind, syscall.Getpid())
	}
	server.ListenAndServe()
}
