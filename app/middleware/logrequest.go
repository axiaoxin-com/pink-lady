package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"

	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}

// LogRequestInfo middleware for logging request info
// set request
// record request and response info
func LogRequestInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ! get request id first
		requestid := c.Request.Header.Get(logging.RequestIDKey)
		if requestid == "" {
			requestid = uuid.NewV4().String()
		}
		c.Set(logging.RequestIDKey, requestid) // 通过上下文环境暴露到handler内部使用
		c.Writer.Header().Set(logging.RequestIDKey, requestid)

		// set ctx logger
		logging.SetCtxLogger(c)

		start := time.Now()
		url := c.Request.URL.String()
		body := ""
		if c.Request.Body != nil {
			buf, _ := ioutil.ReadAll(c.Request.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
			body = readBody(rdr1)
			c.Request.Body = rdr2
		}

		c.Next()

		end := time.Since(start)
		status := c.Writer.Status()
		ctxLogger := logging.CtxLogger(c).With(
			zap.String("url", url),
			zap.String("method", c.Request.Method),
			zap.String("body", body),
			zap.String("clientip", c.ClientIP()),
			zap.String("useragent", c.Request.UserAgent()),
			zap.Int("status", status),
			zap.Int("size", c.Writer.Size()),
			zap.Float64("latency", float64(end.Seconds())*1000.0),
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
