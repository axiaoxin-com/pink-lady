// gin-skeleton: Typically gin-based web application's organizational structure
package main

import (
	"strings"
	"syscall"

	"github.com/axiaoxin/gin-skeleton/app/apis"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	mode := strings.ToLower(viper.GetString("server.mode"))
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()
	apis.RegisterRoutes(app)

	server := endless.NewServer(viper.GetString("server.bind"), app)
	server.BeforeBegin = func(addr string) {
		logrus.Infof("Gin server is listening and serving HTTP on %s (pids: %d)", addr, syscall.Getpid())
	}
	server.ListenAndServe()
}
