package utils

import (
	"net/http"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/sirupsen/logrus"
)

// EndlessServe running a graceful (re)start server
func EndlessServe(bind string, handler http.Handler) {
	server := endless.NewServer(bind, handler)
	server.BeforeBegin = func(bind string) {
		logrus.Infof("Server is listening and serving HTTP on %s (pid: %d)", bind, syscall.Getpid())
	}
	server.ListenAndServe()
}
