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
	r.GET("/", func(c *gin.Context) {})

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
		ctxLogger := logging.GetCtxLogger(c)
		ctxLogger = ctxLogger.With(zap.String("myfield", "myfield"))
		ctxLogger.Info("")
	})

	utils.TestingGETRequest(r, "/")
}

func TestSetNewRequestInfoField(t *testing.T) {
	r := gin.New()
	logging.InitLogger()
	r.Use(LogRequestInfo())
	r.GET("/", func(c *gin.Context) {
		ctxLogger := logging.GetCtxLogger(c)
		ctxLogger.Info("")
		logging.SetCtxLogger(c, zap.String("NewRequestInfoField", "NewRequestInfoField"))
		logging.GetCtxLogger(c).Info("")
	})

	utils.TestingGETRequest(r, "/")
}
