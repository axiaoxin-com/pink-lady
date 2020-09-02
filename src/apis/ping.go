// 默认实现的 ping api

package apis

import (
	"time"

	"github.com/axiaoxin-com/pink-lady/handlers"
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
// @Security BasicAuth
// @Param trace_id header string false "you can set custom trace id in header"
// @Router /x/ping [get]
func Ping(c *gin.Context) {
	time.Sleep(6 * time.Second)
	data := gin.H{
		"version":             Version,
		"mysql":               handlers.CheckMySQL(),
		"redis":               handlers.CheckRedis(),
		"atomic_level_server": handlers.CheckAtomicLevelServer(),
	}
	response.JSON(c, data)
	return
}
