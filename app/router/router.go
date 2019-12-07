// Package router provide a function to setup gin router without routes
package router

import (
	"log"
	"strings"
	"time"

	"pink-lady/app/database"
	"pink-lady/app/logging"
	"pink-lady/app/middleware"
	"pink-lady/app/utils"

	"github.com/spf13/viper"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// InitDependencies 初始化所有依赖
func InitDependencies(configPath, configName string) {
	bindOpt := utils.NewViperOption("server.bind", "localhost:4869", "server binding address")
	modeOpt := utils.NewViperOption("server.mode", "debug", "server mode")
	basicAuthUsernameOpt := utils.NewViperOption("apidocs.basicauth.username", "admin", "apidocs default login username")
	basicAuthPasswordOpt := utils.NewViperOption("apidocs.basicauth.password", "!admin", "apidocs default login password")
	pprofPathOpt := utils.NewViperOption("pprof.path", "/x/debug/pprof", "pprof default register path")
	if err := utils.InitViper(configPath, configName, "",
		bindOpt, modeOpt,
		basicAuthUsernameOpt, basicAuthPasswordOpt,
		pprofPathOpt,
	); err != nil {
		log.Println("[ERROR]", err)
	}

	if err := logging.InitLogger(); err != nil {
		log.Println("[ERROR] ", err)
	}
	if err := utils.InitSentry(); err != nil {
		log.Println("[ERROR] ", err)
	}
	if err := database.InitGorm(); err != nil {
		log.Println("[ERROR] ", err)
	}
	if err := utils.InitRedis(); err != nil {
		log.Println("[ERROR] ", err)
	}
}

// SetupRouter init and return a gin router
func SetupRouter(configPath, configName string) *gin.Engine {
	// Init
	InitDependencies(configPath, configName)

	// setup gin
	mode := strings.ToLower(viper.GetString("server.mode"))
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}
	// new router app
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(middleware.SetRequestID())
	if viper.GetBool("logger.logRequestInfo") {
		router.Use(middleware.LogRequestInfo())
	}
	if viper.GetString("server.sentrydsn") != "" {
		log.Println("[INFO] Using sentry middleware")
		router.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
			Timeout: time.Second * 3,
		}))
	}

	if viper.GetBool("pprof.open") {
		pprof.Register(router, viper.GetString("pprof.path"))
	}

	return router
}
