package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger use the logrus replace gin default logger
func GinLogrus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		c.Next()

		end := time.Since(start)
		status := c.Writer.Status()

		entry := logrus.WithFields(logrus.Fields{
			"path":      path,
			"method":    c.Request.Method,
			"clientIP":  c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
			"requestID": c.MustGet(RequestIDkey),
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
