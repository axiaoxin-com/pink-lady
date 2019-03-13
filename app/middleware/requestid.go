package middleware

import (
	"github.com/axiaoxin/pink-lady/app/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// RequestID is middleware using uuid as request id save in header
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(RequestIDKey)
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}
		c.Set(RequestIDKey, requestID) // 通过上下文环境暴露到handler内部使用

		c.Writer.Header().Set(RequestIDKey, requestID)
		utils.Logger = utils.Logger.WithFields(logrus.Fields{
			"requestID": requestID,
		})
		c.Next()
	}
}
