package main

import (
	"strings"

	"github.com/axiaoxin/gin-skeleton/app/apis"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
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
	endless.ListenAndServe(viper.GetString("server.bind"), app)
}
