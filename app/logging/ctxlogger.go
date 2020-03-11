package logging

import (
	"net/http"
	"net/http/httptest"

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

// GenGinContext 生成gin.context
func GenGinContext(logger *zap.Logger, requestID string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "http://fake.url", nil)
	SetCtxLogger(c, logger)
	SetCtxRequestID(c, requestID)
	return c
}

// SetCtxLogger add a logger with given field into the gin.Context
// and set requestid field get from context
func SetCtxLogger(c *gin.Context, ctxLogger *zap.Logger) {
	c.Set(CtxLoggerKey, ctxLogger)
}

// CtxLogger get the ctxLogger in gin.Context
func CtxLogger(c *gin.Context, fields ...zap.Field) *zap.Logger {
	var ctxLogger *zap.Logger
	ctxLoggerItf, _ := c.Get(CtxLoggerKey)
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
func CtxRequestID(c *gin.Context) string {
	// first get from context
	requestid := c.GetString(RequestIDKey)
	if requestid != "" {
		return requestid
	}
	Debug("no requestid in context")
	// if not then get request id from header
	requestid = c.Request.Header.Get(RequestIDKey)
	if requestid != "" {
		return requestid
	}
	Debug("no requestid in header")
	// else gen a request id
	return "pink-lady-" + uuid.NewV4().String()
}

// SetCtxRequestID 设置requestid到ctxlogger、context、header中
func SetCtxRequestID(c *gin.Context, requestid string) {
	// 设置带requestid的logger到context中
	ctxLogger := CtxLogger(c, zap.String(RequestIDKey, requestid))
	SetCtxLogger(c, ctxLogger)
	// 设置requestid到context中
	c.Set(RequestIDKey, requestid)
	// 设置requestid到header中
	c.Request.Header.Set(RequestIDKey, requestid)
	c.Writer.Header().Set(RequestIDKey, requestid)
}
