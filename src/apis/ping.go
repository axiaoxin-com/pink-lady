// 默认实现的 ping api

package apis

import (
	"github.com/axiaoxin-com/pink-lady/response"

	"github.com/gin-gonic/gin"
)

// VERSION your API version, don't change the code style
const VERSION = "0.0.1"

// Ping godoc
// @Summary Ping for server is living will respond API version
// @Tags x
// @Produce json
// @Router /x/ping [get]
// @Success 200 {object} response.Response
func Ping(c *gin.Context) {
	data := gin.H{"version": VERSION}
	response.JSON(c, data)
}
