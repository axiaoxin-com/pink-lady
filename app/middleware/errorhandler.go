package middleware

import (
	"net/http"

	"gin-skeleton/app/services/retcode"
	"gin-skeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		switch status {
		case http.StatusNotFound:
			response.JSON404(c, retcode.RouteNotFound, c.Errors.String())
		case http.StatusInternalServerError:
			response.JSON500(c, retcode.InternalError, c.Errors.String())
		default:
			if status > http.StatusBadRequest {
				response.Respond(c, status, retcode.UnknownError, c.Errors.String())
			}

		}
	}
}
