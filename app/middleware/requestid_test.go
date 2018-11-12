package middleware

import (
	"testing"

	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	r := gin.New()
	r.Use(RequestID())
	r.GET("/", func(c *gin.Context) {})
	w := utils.PerformTestingRequest(r, "GET", "/")
	if w.Header().Get(RequestIDKey) == "" {
		t.Error("no request id")
	}
}
