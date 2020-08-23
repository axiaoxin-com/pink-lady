package main

import (
	"github.com/axiaoxin-com/pink-lady/apis"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
)

func middlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{}
	return m
}

func main() {
	app := webserver.NewGinEngine(middlewares()...)
	webserver.Run(app, apis.Register)
}
