// Package router provide a function to setup gin router without routes
package router

import (
	"strings"
	"time"

	"github.com/axiaoxin/pink-lady/app/middleware"
	"github.com/spf13/viper"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter init and return a gin router
func SetupRouter() *gin.Engine {
	mode := strings.ToLower(viper.GetString("server.mode"))

	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(middleware.LogRequestInfo())
	if sentry.CurrentHub().Client() != nil {
		router.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
			Timeout: time.Second * 3,
		}))
	}
	return router
}
