package middleware

import (
	"testing"

	"github.com/axiaoxin/pink-lady/app/utils"

	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	r := gin.New()
	r.Use(RequestID())
	r.GET("/", func(c *gin.Context) {})
	w := utils.TestingGETRequest(r, "/")
	if w.Header().Get(RequestIDKey) == "" {
		t.Error("no request id")
	}
}
