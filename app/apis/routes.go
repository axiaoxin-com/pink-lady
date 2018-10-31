// package apis routes.go provides register handlers on url
package apis

import (
	_ "github.com/axiaoxin/gin-skeleton/app/apis/docs"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// RegisterRoutes add handlers on urls at there
func RegisterRoutes(app *gin.Engine) {

	// group x registered gin-skeleton default api
	x := app.Group("/x")
	{
		x.GET("/apidocs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		x.GET("/ping", Ping)
	}

	// register your api below
}
