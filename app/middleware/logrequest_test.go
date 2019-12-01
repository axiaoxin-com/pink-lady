package middleware

import (
	"testing"

	"pink-lady/app/logging"
	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
)

func TestLogRequestInfo(t *testing.T) {
	r := gin.New()
	logging.InitLogger()
	r.Use(LogRequestInfo())
	r.GET("/", func(c *gin.Context) { return })
	utils.PerformRequest(r, "GET", "/", []byte("request\t\nbody "))
	// show log request info fields
}

func Test500(t *testing.T) {
	r := gin.New()
	logging.InitLogger()
	r.Use(LogRequestInfo())
	r.GET("/200", func(c *gin.Context) {
		return
	})
	r.GET("/500", func(c *gin.Context) {
		c.AbortWithStatus(500)
		return
	})

	utils.PerformRequest(r, "GET", "/200", nil)
	utils.PerformRequest(r, "GET", "/500", nil)
	// show error level log
}
