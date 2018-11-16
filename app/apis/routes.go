// Package apis routes.go provides register handlers on url
package apis

import (
	"pink-lady/app/apis/demo"
	// need by swag
	_ "pink-lady/app/docs"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
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

	// group demo registered a demo, you can delete it
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

	// register your api below
}
