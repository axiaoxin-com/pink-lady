package logging

import (
	"net/http"
	"net/http/httptest"
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

	SetCtxLogger(c, Logger)
	logger := CtxLogger(c)
	if logger == nil {
		t.Fatal("context also must should return a logger")
	}
	if logger != Logger {
		t.Fatal("logger no equal")
	}
	logger.Info("this is a logger from default logger")
}

func TestCtxLoggerDefaultLoggerWithField(t *testing.T) {
	InitLogger()
	c := &gin.Context{}

	SetCtxLogger(c, Logger)
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

func TestCtxRequstIDHeader(t *testing.T) {
	InitLogger()
	c := &gin.Context{}
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set(RequestIDKey, "IAMAREQUESTID TOO")
	if CtxRequestID(c) != "IAMAREQUESTID TOO" {
		t.Fatal("context should return set value")
	}
}

func TestSetRequestID(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/", nil)
	SetCtxRequestID(c, "xyz")
	r, e := c.Get(RequestIDKey)
	if !e {
		t.Fatal("xyz not exists")
	}
	if r.(string) != "xyz" {
		t.Fatal("should xyz")
	}
	if c.Request.Header.Get(RequestIDKey) != "xyz" {
		t.Fatal("request header not xyz")
	}
	if c.Writer.Header().Get(RequestIDKey) != "xyz" {
		t.Fatal("request header not xyz")
	}
}
