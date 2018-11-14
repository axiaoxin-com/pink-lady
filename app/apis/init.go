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
	"github.com/axiaoxin/gin-skeleton/app/services/retcode"
	"github.com/axiaoxin/gin-skeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

// API version
const VERSION = "0.0.1"

// package init function
func init() {
}

// Ping godoc
// @Summary Ping for server is living
// @Description response current api version
// @Tags x
// @Produce  json
// @Router /x/ping [get]
// @Success 200 {object} response.Response
func Ping(c *gin.Context) {
	data := gin.H{"version": VERSION}
	response.JSON(c, retcode.Success, data)
}
