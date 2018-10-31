package middleware

import (
	"net/http"

	"github.com/axiaoxin/gin-skeleton/app/apis"
	"github.com/axiaoxin/gin-skeleton/app/common/retcode"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		switch status {
		case http.StatusNotFound:
			apis.Respond(c, retcode.APINotFound, c.Errors.String())
		case http.StatusInternalServerError:
			apis.Respond(c, retcode.InternalError, c.Errors.String())
		default:
			if status > http.StatusBadRequest {
				apis.Respond(c, retcode.UnknownError, c.Errors.String())
			}

		}
	}
}
