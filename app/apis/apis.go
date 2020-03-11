// Package apis save the web api code
// init.go define init() function and default system api
// handlers save you handle functions
// routes.go register handle function on url
// WAY TO ADD YOUR NEW API:
// create code file or package according to you business logic, let structure be modularized
// write the gin handlerFunc code like the Ping() in the file
// you should extract the common business logic handle functions into handlers package
// database model should be defined in models package by modularized
// general tool functions should be defined in utils package by modularized
// you can return unified json struct by response package
// the new return code should be defined in retcode package
// when you finish the handlerFunc you need to register it on a url in routes.go
// that's all.
package apis

import (
	"pink-lady/app/logging"
	"pink-lady/app/response"

	"github.com/gin-gonic/gin"
)

// @title pink-lady Web API
// @version 0.0.1
// @description pink-lady web API list.
// @contact.name API Support
// @contact.url http://axiaoxin.com
// @contact.email 254606826@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:4869
// @BasePath /

// VERSION your API version, don't change the code style
const VERSION = "0.0.1"

// Ping godoc
// @Summary Ping for server is living will respond API version
// @Tags x
// @Produce  json
// @Router /x/ping [get]
// @Success 200 {object} response.Response
func Ping(c *gin.Context) {
	data := gin.H{"version": VERSION}
	logging.CtxLogger(c).Debug("===> Ping")
	response.JSON(c, data)
}
