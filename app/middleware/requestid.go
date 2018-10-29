package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "X-Request-Id"
		requestID := c.Request.Header.Get(key)
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}
		c.Set(key, requestID) // 暴露到handler内部使用

		c.Writer.Header().Set(key, requestID)
		c.Next()
	}
}
