package logging

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetCtxLogger(t *testing.T) {
	InitLogger()
	c := &gin.Context{}
	SetCtxLogger(c)
	_, exists := c.Get(CtxLoggerKey)
	if !exists {
		t.Error("ctxLogger not exsits")
	}
	logger := CtxLogger(c)
	if logger == nil {
		t.Error("ctxLogger is nil")
	}
}
