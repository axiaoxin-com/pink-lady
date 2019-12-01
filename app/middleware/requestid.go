package middleware

import (
	"pink-lady/app/logging"

	"github.com/gin-gonic/gin"
)

// SetRequestID middleware for gen request id
func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从context获取requestid，不存在则新生成
		requestid := logging.CtxRequestID(c)
		// 设置requestid到ctxlogger、context、header中
		logging.SetCtxRequestID(c, requestid)

		c.Next()
	}
}
