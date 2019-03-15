package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/axiaoxin/pink-lady/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}

// GinLogrus is a logger middleware, which use the logrus replace gin default logger
func GinLogrus() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		entry := utils.Logger.WithFields(logrus.Fields{
			"url":       url,
			"method":    c.Request.Method,
			"body":      body,
			"clientIP":  c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
			"status":    status,
			"size":      c.Writer.Size(),
			"latency":   fmt.Sprintf("%fms", float64(end.Seconds())*1000.0),
		})

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			if status > 499 {
				entry.Error()
			} else {
				entry.Info()
			}
		}
	}
}
