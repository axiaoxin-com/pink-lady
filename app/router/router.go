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

const (
	// DefaultPprofPath 默认pprof url
	DefaultPprofPath = "/x/debug/pprof"
)

// InitDependencies 初始化所有依赖
func InitDependencies(configpath, configname string) {
	if err := utils.InitViper(configpath, configname, ""); err != nil {
		log.Println("[ERROR]", err)
	}

	if err := logging.InitLogger(); err != nil {
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
func SetupRouter(configpath, configname string) *gin.Engine {
	// Init
	InitDependencies(configpath, configname)

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
		pprofPath := viper.GetString("pprof.path")
		if pprofPath == "" {
			pprofPath = DefaultPprofPath
		}
		pprof.Register(router, pprofPath)
	}

	return router
}
