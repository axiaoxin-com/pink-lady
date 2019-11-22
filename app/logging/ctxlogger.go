package logging

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// CtxLoggerKey define the logger keyname which in context
	CtxLoggerKey = "ctxLogger"
	// RequestIDKey define the request id header key
	RequestIDKey = "X-Request-Id"
)

// SetCtxLogger add a logger with given field into the gin.Context
// and set requestid field get from context
func SetCtxLogger(c *gin.Context, fields ...zap.Field) {
	ctxLogger := CloneLogger()
	ctxLogger = ctxLogger.With(fields...)
	requestid, _ := c.Get(RequestIDKey)
	if requestid != nil {
		ctxLogger = ctxLogger.With(zap.String("requestid", requestid.(string)))
	}
	c.Set(CtxLoggerKey, ctxLogger)
}

// CtxLogger get the ctxLogger in gin.Context
func CtxLogger(c *gin.Context) *zap.Logger {
	ctxLogger, _ := c.Get(CtxLoggerKey)
	return ctxLogger.(*zap.Logger)
}
