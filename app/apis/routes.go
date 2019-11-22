// Package apis routes.go provides register handlers on url
package apis

import (
	"github.com/axiaoxin/pink-lady/app/apis/demo"
	// need by swag
	_ "github.com/axiaoxin/pink-lady/app/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// RegisterRoutes add handlers on urls at there
func RegisterRoutes(app *gin.Engine) {

	// group x registered pink-lady default api
	x := app.Group("/x")
	{
		x.GET("/apidocs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		x.GET("/ping", Ping)
	}

	/*
		// redirect / to apidocs
		app.GET("/", func(c *gin.Context) {
			c.Redirect(301, "/x/apidocs/index.html")
		})
	*/

	// demo routes start
	demoGroup := app.Group("/demo")
	{
		demoGroup.POST("/label", demo.AddLabel)
		demoGroup.GET("/label", demo.Label)

		demoGroup.POST("/object", demo.AddObject)
		demoGroup.GET("/object", demo.Object)

		demoGroup.POST("/labeling", demo.AddLabeling)
		demoGroup.GET("/labeling/label/:id", demo.GetLabelingByLabelID)
		demoGroup.GET("/labeling/object/:id", demo.GetLabelingByObjectID)
		demoGroup.PUT("/labeling", demo.ReplaceLabeling)
		demoGroup.DELETE("/labeling", demo.DeleteLabeling)
	}
	// demo routes end
	// register your api below
}
