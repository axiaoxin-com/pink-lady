// package web routes.go provides register handlers on url
package main

import (
	"github.com/axiaoxin/gin-skeleton/app/apis"
	"github.com/gin-gonic/gin"
)

func registerRoutes(app *gin.Engine) {
	app.GET("/ping", apis.Ping)
}
