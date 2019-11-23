package logging

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestSetCtxLogger(t *testing.T) {
	InitLogger()
	c := &gin.Context{}

	SetCtxLogger(c, Logger)
	_, exists := c.Get(CtxLoggerKey)
	if !exists {
		t.Error("set ctxLogger failed")
	}
}

func TestCtxLoggerEmpty(t *testing.T) {
	InitLogger()
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)

	logger := CtxLogger(c)
	if logger == nil {
		t.Error("empty context also must should return a logger")
	}
	logger.Info("this is a logger from empty ctx")
}

func TestCtxLoggerEmptyField(t *testing.T) {
	InitLogger()
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)

	logger := CtxLogger(c, zap.String("field1", "1"))
	if logger == nil {
		t.Error("empty context also must should return a logger")
	}
	logger.Info("this is a logger from empty ctx but with field")
}

func TestCtxLoggerDefaultLogger(t *testing.T) {
	InitLogger()
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)

	SetCtxLogger(c, Logger)
	logger := CtxLogger(c)
	if logger == nil {
		t.Error("context also must should return a logger")
	}
	logger.Info("this is a logger from default logger")
}

func TestCtxLoggerDefaultLoggerWithField(t *testing.T) {
	InitLogger()
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)

	SetCtxLogger(c, Logger)
	logger := CtxLogger(c, zap.String("myfield", "xxx"))
	if logger == nil {
		t.Error("context also must should return a logger")
	}
	logger.Info("this is a logger from default logger with field")
}

func TestCtxLoggerWithCtxReuquestID(t *testing.T) {
	InitLogger()
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Set(RequestIDKey, "IAMAREQUESTID")

	logger := CtxLogger(c)
	if logger == nil {
		t.Error("context also must should return a logger")
	}
	logger.Info("this is a logger context exists requestid")
}

func TestCtxLoggerWithHeaderReuquestID(t *testing.T) {
	InitLogger()
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set(RequestIDKey, "IAMAREQUESTID TOO")

	logger := CtxLogger(c)
	if logger == nil {
		t.Error("context also must should return a logger")
	}
	logger.Info("this is a logger header exists requestid")
}
