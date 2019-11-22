package utils

import (
	"log"
	"net/http"
	"syscall"

	"github.com/fvbock/endless"
)

// EndlessServe running a graceful (re)start server
func EndlessServe(bind string, handler http.Handler) {
	server := endless.NewServer(bind, handler)
	server.BeforeBegin = func(bind string) {
		log.Printf("Server is listening and serving HTTP on %s (pid: %d)\n", bind, syscall.Getpid())
	}
	server.ListenAndServe()
}
