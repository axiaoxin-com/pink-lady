package middleware

import (
	"testing"

	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/axiaoxin/pink-lady/app/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	r := gin.New()
	logging.InitLogger()
	r.Use(LogRequestInfo())
	r.GET("/", func(c *gin.Context) { return })

	w := utils.TestingGETRequest(r, "/")
	if w.Header().Get(logging.RequestIDKey) == "" {
		t.Error("no request id")
	}
}

func TestCtxLogger(t *testing.T) {
	r := gin.New()
	logging.InitLogger()
	r.Use(LogRequestInfo())
	r.GET("/", func(c *gin.Context) {
		ctxLogger := logging.CtxLogger(c)
		ctxLogger = ctxLogger.With(zap.String("myfield", "myfield"))
		ctxLogger.Info("")
		return
	})

	utils.TestingGETRequest(r, "/")
}

func TestSetNewRequestInfoField(t *testing.T) {
	r := gin.New()
	logging.InitLogger()
	r.Use(LogRequestInfo())
	r.GET("/", func(c *gin.Context) {
		ctxLogger := logging.CtxLogger(c)
		ctxLogger.Info("")
		logging.SetCtxLogger(c, zap.String("NewRequestInfoField", "NewRequestInfoField"))
		logging.CtxLogger(c).Info("")
		return
	})

	utils.TestingGETRequest(r, "/")
}

func Test500(t *testing.T) {
	r := gin.New()
	logging.InitLogger()
	r.Use(LogRequestInfo())
	r.GET("/200", func(c *gin.Context) {
		t.Log("200")
		return
	})
	r.GET("/500", func(c *gin.Context) {
		c.AbortWithStatus(500)
		return
	})

	utils.TestingGETRequest(r, "/500")
	utils.TestingGETRequest(r, "/200")
}
