// Package router provide a function to setup gin router without routes
package router

import (
	"strings"

	"github.com/axiaoxin/pink-lady/app/middleware"
	"github.com/spf13/viper"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

// SetupRouter init and return a gin router
func SetupRouter() *gin.Engine {
	mode := strings.ToLower(viper.GetString("server.mode"))
	sentryDSN := viper.GetString("sentry.dsn")
	sentryOnlyCrashes := viper.GetBool("sentry.onlycrashes")
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(cors.Default())
	router.Use(middleware.LogRequestInfo())
	if sentryDSN != "" {
		raven.SetDSN(sentryDSN)
		router.Use(sentry.Recovery(raven.DefaultClient, sentryOnlyCrashes))
	}
	return router
}
