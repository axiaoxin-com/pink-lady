// package apis routes.go provides register handlers on url
package apis

import (
	_ "github.com/axiaoxin/gin-skeleton/app/apis/docs"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func RegisterRoutes(app *gin.Engine) {
	// docs url
	docs := app.Group("/docs")
	{
		docs.GET("/apis/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// group x registered gin-skeleton default api
	x := app.Group("/x")
	{
		x.GET("/ping", Ping)
	}

	// register your api below
}
