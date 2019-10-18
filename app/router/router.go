// Package router provide a function to setup gin router without routes
package router

import (
	"github.com/axiaoxin/pink-lady/app/middleware"

	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

// SetupRouter init and return a gin router
func SetupRouter(mode string, sentryDSN string, sentryOnlyCrashes bool) *gin.Engine {
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
	router.Use(middleware.RequestID()) // requestid 必须在ginlogrus前面
	router.Use(middleware.GinLogrus())
	if sentryDSN != "" {
		raven.SetDSN(sentryDSN)
		router.Use(sentry.Recovery(raven.DefaultClient, sentryOnlyCrashes))
	}
	return router
}
