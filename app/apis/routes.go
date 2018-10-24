// package apis routes.go provides register handlers on url
package apis

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *gin.Engine) {
	app.GET("/ping", Ping)
}
