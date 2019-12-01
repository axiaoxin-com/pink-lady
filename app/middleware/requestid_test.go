package middleware

import (
	"testing"

	"pink-lady/app/logging"
	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	r := gin.New()
	logging.InitLogger()
	r.Use(SetRequestID())
	r.GET("/", func(c *gin.Context) {
		_, e := c.Get(logging.RequestIDKey)
		if !e {
			t.Fatal("Context中没有requestid")
		}
		return
	})

	respRecorder := utils.PerformRequest(r, "GET", "/", nil)
	if respRecorder.Header().Get(logging.RequestIDKey) == "" {
		t.Fatal("no request id")
	}
}
