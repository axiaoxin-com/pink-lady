// package apis save the web api code
// init.go define init() function and default system api
// handlers save you handle functions
// routes.go register handle function on url
//
// WAY TO ADD YOUR NEW API:
// create code file or package according to you business logic, let structure be modularized
// write the gin handlerFunc code like the Ping() in the file
// you should extract the common business logic handle functions into services package
// database model should be defined in models package by modularized
// general tool functions should be defined in utils package by modularized
// you can record log by logrus and get config by viper
// you can return unified json struct by utils/response package
// the new return code should be defined in services/retcode package
// when you finish the handlerFunc you need to register it on a url in routes.go
// that's all.

// @title Gin-Skeleton Web API
// @version 0.0.1
// @description These are web APIs based on gin-skeleton.

// @contact.name API Support
// @contact.url http://km.oa.com/user/ashinchen
// @contact.email ashinchen@tencent.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
package apis

import (
	"github.com/axiaoxin/gin-skeleton/app/middleware"
	"github.com/axiaoxin/gin-skeleton/app/services"
	"github.com/axiaoxin/gin-skeleton/app/services/retcode"
	"github.com/axiaoxin/gin-skeleton/app/utils/response"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

// package init function
func init() {

}

// SetupRouter init gin register routes and return a router app
func SetupRouter(mode string, sentryDSN string, sentryOnlyCrashes bool) *gin.Engine {
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.Use(middleware.ErrorHandler())
	app.Use(cors.Default())
	app.Use(middleware.RequestID())
	app.Use(middleware.GinLogrus())
	if sentryDSN != "" {
		raven.SetDSN(sentryDSN)
		app.Use(sentry.Recovery(raven.DefaultClient, sentryOnlyCrashes))
	}

	RegisterRoutes(app)
	return app
}

// Ping godoc
// @Summary Ping for server is living
// @Description response current api version
// @Tags x
// @Produce  json
// @Router /x/ping [get]
// @Success 200 {object} response.Response
func Ping(c *gin.Context) {
	data := gin.H{"version": services.VERSION}
	response.JSON(c, retcode.Success, data)
}
