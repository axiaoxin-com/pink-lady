package main

import (
	"strings"

	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func initAPP() *gin.Engine {
	initViper()

	utils.InitLogrus(viper.GetString("log_level"), viper.GetString("log_formatter"), viper.GetString("log_out"))

	// init gin
	mode := strings.ToLower(viper.GetString("mode"))
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// init api server
	// app := gin.New()
	//app.Use(logger, gin.Recovery())
	app := gin.Default()
	registerRoutes(app)
	return app
}
