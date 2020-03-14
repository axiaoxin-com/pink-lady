package logging

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestSetCtxLoggerRequestID(t *testing.T) {
	InitLogger()
	c := &gin.Context{}

	SetCtxLogger(c, "1234")
	_, exists := c.Get(CtxLoggerKey)
	if !exists {
		t.Fatal("set ctxLogger failed")
	}
}

func TestCtxLoggerEmpty(t *testing.T) {
	InitLogger()
	c := &gin.Context{}

	logger := CtxLogger(c)
	if logger == nil {
		t.Fatal("empty context also must should return a logger")
	}
	logger.Info("this is a logger from empty ctx")
}

func TestCtxLoggerEmptyField(t *testing.T) {
	InitLogger()
	c := &gin.Context{}

	logger := CtxLogger(c, zap.String("field1", "1"))
	if logger == nil {
		t.Fatal("empty context also must should return a logger")
	}
	logger.Info("this is a logger from empty ctx but with field")
}

func TestCtxLoggerDefaultLogger(t *testing.T) {
	InitLogger()
	c := &gin.Context{}

	SetCtxLogger(c, "rid")
	logger := CtxLogger(c)
	if logger == nil {
		t.Fatal("context also must should return a logger")
	}
	logger.Info("this is a logger from default logger")
}

func TestCtxLoggerDefaultLoggerWithField(t *testing.T) {
	InitLogger()
	c := &gin.Context{}

	SetCtxLogger(c, "rid")
	logger := CtxLogger(c, zap.String("myfield", "xxx"))
	if logger == nil {
		t.Fatal("context also must should return a logger")
	}
	if logger == Logger {
		t.Fatal("with field will get a logger")
	}
	logger.Info("this is a logger from default logger with field")
}

func TestCtxRequstIDCtx(t *testing.T) {
	InitLogger()
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)
	if CtxRequestID(c) == "" {
		t.Fatal("context should return default value")
	}
	c.Set(RequestIDKey, "IAMAREQUESTID")
	if CtxRequestID(c) != "IAMAREQUESTID" {
		t.Fatal("context should return set value")
	}
}
