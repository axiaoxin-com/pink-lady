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
			apis.JSON404(c, retcode.RouteNotFound, c.Errors.String())
		case http.StatusInternalServerError:
			apis.JSON500(c, retcode.InternalError, c.Errors.String())
		default:
			if status > http.StatusBadRequest {
				apis.Respond(c, status, retcode.UnknownError, c.Errors.String())
			}

		}
	}
}
