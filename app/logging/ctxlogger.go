package logging

import (
	"context"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

const (
	// CtxLoggerKey define the logger keyname which in context
	CtxLoggerKey = "ctxLogger"
	// RequestIDKey define the request id header key
	RequestIDKey = "X-Request-Id"
)

// CtxLogger get the ctxLogger in context
func CtxLogger(c context.Context, fields ...zap.Field) *zap.Logger {
	if c == nil {
		c = context.Background()
	}
	var ctxLogger *zap.Logger
	var ctxLoggerItf interface{}
	if gc, ok := c.(*gin.Context); ok {
		ctxLoggerItf, _ = gc.Get(CtxLoggerKey)
	} else {
		ctxLoggerItf = c.Value(CtxLoggerKey)
	}

	if ctxLoggerItf != nil {
		ctxLogger = ctxLoggerItf.(*zap.Logger)
	} else {
		ctxLogger = CloneLogger().Named("ctxLogger")
		Debug("no ctxLogger in context, clone the global Logger as ctxLogger")
	}
	if len(fields) > 0 {
		ctxLogger = ctxLogger.With(fields...)
	}
	return ctxLogger
}

// CtxRequestID get requestid from context
func CtxRequestID(c context.Context) string {
	// first get from gin context
	if gc, ok := c.(*gin.Context); ok {
		if requestid := gc.GetString(RequestIDKey); requestid != "" {
			return requestid
		}
	}
	// get from go context
	requestidItf := c.Value(RequestIDKey)
	if requestidItf != nil {
		return requestidItf.(string)
	}
	// return default value
	return "pink-lady-" + uuid.NewV4().String()
}

// SetCtxLogger set logger and requestid into context
func SetCtxLogger(c context.Context, requestid string) {
	logger := CtxLogger(c, zap.String(RequestIDKey, requestid))
	if gc, ok := c.(*gin.Context); ok {
		gc.Set(CtxLoggerKey, logger)
	}
	c = context.WithValue(c, CtxLoggerKey, logger)

}
