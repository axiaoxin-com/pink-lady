package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"pink-lady/app/logging"
	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogRequestInfo middleware for logging request info
// set request
// record request and response info
func LogRequestInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 读取body用于日志打印
		var body []byte
		if c.Request.Body != nil {
			body, _ = ioutil.ReadAll(c.Request.Body)
		}
		// body被read、bind之后会被置空，需要重置
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		c.Next()

		// 执行完毕后打印详情
		end := time.Since(start)
		status := c.Writer.Status()
		ctxLogger := logging.CtxLogger(c,
			zap.String("url", c.Request.URL.String()),
			zap.String("method", c.Request.Method),
			zap.Int("status", status),
			zap.Float64("latency", float64(end.Seconds())*1000.0), // 毫秒
			zap.String("clientip", c.ClientIP()),
			zap.String("useragent", c.Request.UserAgent()),
			zap.Int("size", c.Writer.Size()),
			zap.String("body", utils.RemoveAllWhiteSpace(string(body))),
		)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			ctxLogger.Error(c.Errors.String())
		} else {
			if status > 499 {
				ctxLogger.Error("LogRequestInfo")
			} else {
				ctxLogger.Info("LogRequestInfo")
			}
		}
	}
}
