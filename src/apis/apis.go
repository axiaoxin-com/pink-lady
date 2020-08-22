// @title apidocs
// @version 0.0.1
// @description web apis written with pink-lady
// @contact.name API Support
// @contact.url http://axiaoxin.com
// @contact.email 254606826@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:4869
// @BasePath /

package apis

import (
	// need by swag
	_ "github.com/axiaoxin-com/pink-lady/apis/docs"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Register 在 gin engine 上注册 url 对应的 HandlerFunc
func Register(app *gin.Engine) {
	// auth 加到路由中可以对该路由添加 basic auth 登录验证
	auth := gin.BasicAuth(gin.Accounts{
		viper.GetString("basic_auth.username"): viper.GetString("basic_auth.password"),
	})
	// Group x 添加了一些默认的 url 路由
	x := app.Group("/x")
	{
		// ginSwagger 生成的在线 API 文档路由
		x.GET("/apidocs/*any", auth, ginSwagger.WrapHandler(swaggerFiles.Handler))
		// 默认的 ping 方法，返回 server 相关信息
		x.GET("/ping", Ping)
	}

	// 在这下面开始添加你的 gin HandlerFunc
}
