// package apis routes.go provides register handlers on url
package apis

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *gin.Engine) {
	// group x registered gin-skeleton default api
	x := app.Group("/x")
	{
		x.GET("/ping", Ping)
	}

	// register your api below
}
