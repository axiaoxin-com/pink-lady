// Package apis routes.go provides register handlers on url
package apis

import (
	// need by swag
	_ "pink-lady/app/apis/docs"
	"pink-lady/app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// RegisterRoutes add handlers on urls at there
func RegisterRoutes(app *gin.Engine) {
	// group x registered pink-lady default api
	x := app.Group("/x")
	{
		x.GET("/apidocs/*any", middleware.BasicAuth(), ginSwagger.WrapHandler(swaggerFiles.Handler))
		x.GET("/ping", Ping)
	}

	// redirect / to apidocs
	if viper.GetBool("apidocs.rootRedirect") {
		app.GET("/", func(c *gin.Context) {
			c.Redirect(301, "/x/apidocs/index.html")
		})
	}

	// demo routes start
	demoGroup := app.Group("/demo")
	{
		demoGroup.POST("/alert-policy", CreateAlertPolicy)
		demoGroup.GET("/alert-policy", DescribeAlertPolicies)
		demoGroup.PUT("/alert-policy", ModifyAlertPolicy)
		demoGroup.GET("/alert-policy/:appid/:uin/:id", DescribeAlertPolicy)
		demoGroup.DELETE("/alert-policy/:appid/:uin/:id", DeleteAlertPolicy)
	}
	// demo routes end

	// 在这下面开始注册你的URL路由
}
