package response

import (
	"net/http"

	"github.com/axiaoxin/gin-skeleton/app/services/retcode"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JSON respond unified JSON structure with 200 http status code
func JSON(c *gin.Context, rc *retcode.RetCode, data interface{}) {
	Respond(c, http.StatusOK, rc, data)
}

// JSON400 respond unified JSON structure with 400 http status code
func JSON400(c *gin.Context, rc *retcode.RetCode, data interface{}) {
	Respond(c, http.StatusBadRequest, rc, data)
}

// JSON404 respond unified JSON structure with 404 http status code
func JSON404(c *gin.Context, rc *retcode.RetCode, data interface{}) {
	Respond(c, http.StatusNotFound, rc, data)
}

// JSON500 respond unified JSON structure with 500 http status code
func JSON500(c *gin.Context, rc *retcode.RetCode, data interface{}) {
	Respond(c, http.StatusInternalServerError, rc, data)
}

// Respond encapsulates c.JSON
// debug mode respond indented json
func Respond(c *gin.Context, status int, rc *retcode.RetCode, data interface{}) {
	code, msg := rc.Decode()
	resp := Response{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	if gin.Mode() == gin.ReleaseMode {
		c.JSON(status, resp)
	} else {
		c.IndentedJSON(status, resp)
	}
}
