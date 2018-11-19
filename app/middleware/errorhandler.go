package middleware

import (
	"net/http"

	"github.com/axiaoxin/pink-lady/app/services/retcode"
	"github.com/axiaoxin/pink-lady/app/utils/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler return json on 404 or 500
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
