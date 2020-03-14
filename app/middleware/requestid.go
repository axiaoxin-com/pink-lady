package middleware

import (
	"pink-lady/app/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetRequestID middleware for gen request id
func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从context获取requestid，不存在则新生成
		requestid := c.Request.Header.Get(logging.RequestIDKey)
		if requestid == "" {
			requestid = logging.CtxRequestID(c)
		}
		// 设置requestid到header中
		c.Writer.Header().Set(logging.RequestIDKey, requestid)
		// 设置requestid到gin context中
		c.Set(logging.RequestIDKey, requestid)
		// 设置requestid到ctxlogger中
		ctxLogger := logging.CtxLogger(c, zap.String(logging.RequestIDKey, requestid))
		// 将logger设置到gin context
		logging.SetCtxLogger(c, ctxLogger)

		c.Next()
	}
}
