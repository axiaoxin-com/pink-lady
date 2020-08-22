package main

import (
	"github.com/axiaoxin-com/pink-lady/apis"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	webserver.InitConfig()
}

func middlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{}
	return m
}

func main() {
	app := webserver.NewGinEngine(viper.GetString("server.mode"), viper.GetBool("server.pprof"), middlewares()...)
	webserver.Run(app, apis.Register)
}
