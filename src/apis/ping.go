// 默认实现的 ping api

package apis

import (
	"github.com/axiaoxin-com/pink-lady/response"

	"github.com/gin-gonic/gin"
)

// Ping godoc
// @Summary 默认的 Ping 接口
// @Description 返回 server 相关信息，可以用于健康检查
// @Tags x
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
// @Router /x/ping [get]
func Ping(c *gin.Context) {
	data := gin.H{"version": Version}
	response.JSON(c, data)
	panic("test panic")
}
