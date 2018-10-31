package apis

import (
	"net/http"

	"github.com/axiaoxin/gin-skeleton/app/common/retcode"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Respond encapsulates c.JSON
// debug mode respond indented json
func Respond(c *gin.Context, rc *retcode.RetCode, data interface{}) {
	code, msg := rc.Decode()
	if gin.Mode() == gin.ReleaseMode {
		c.JSON(http.StatusOK, Response{
			Code:    code,
			Message: msg,
			Data:    data,
		})
	} else {
		c.IndentedJSON(http.StatusOK, Response{
			Code:    code,
			Message: msg,
			Data:    data,
		})
	}
}
