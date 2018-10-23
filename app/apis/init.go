// skeleton system default apis in here
package apis

import (
	"net/http"

	"github.com/axiaoxin/gin-skeleton/app/handlers"
	"github.com/gin-gonic/gin"
)

// response current api version for ping request
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":    handlers.VERSION,
		"message": "success",
		"code":    0,
	})
}
