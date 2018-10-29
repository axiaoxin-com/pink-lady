package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// use uuid as request id
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(RequestIDkey)
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}
		c.Set(RequestIDkey, requestID) // 通过上下文环境暴露到handler内部使用

		c.Writer.Header().Set(RequestIDkey, requestID)
		c.Next()
	}
}
