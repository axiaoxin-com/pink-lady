package middleware

import (
	"testing"

	"go.uber.org/zap"

	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/axiaoxin/pink-lady/app/utils"
	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	r := gin.New()
	logging.InitLogger()
	r.Use(LogRequestInfo())
	r.GET("/", func(c *gin.Context) {
		_, e := c.Get(logging.RequestIDKey)
		if !e {
			t.Error("Context中没有requestid")
		}
		return
	})

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
		ctxLogger.Info("from ctxlogger var")
		logging.CtxLogger(c).Info("from middleware logger")
		ctxLogger = ctxLogger.With(zap.String("myfield", "myfield"))
		ctxLogger.Info("from ctxlogger var with field")
		logging.CtxLogger(c).Info("from middleware logger after var set field")

		dctxLogger := logging.CtxLogger(c)
		nctxLogger := dctxLogger.With(zap.String("new-field", "new_field"))
		nctxLogger.Info("from nctxLogger with field")
		logging.SetCtxLogger(c, nctxLogger)
		logging.CtxLogger(c).Info("from middleware logger after set new logger")
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
