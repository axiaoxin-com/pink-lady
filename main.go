//go:generate swag init --dir ./ --generalInfo routes/routes.go --propertyStrategy snakecase --output ./routes/docs

package main

import (
	"flag"

	"github.com/axiaoxin-com/pink-lady/models"
	"github.com/axiaoxin-com/pink-lady/routes"
	"github.com/axiaoxin-com/pink-lady/routes/response"
	"github.com/axiaoxin-com/pink-lady/services"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// DefaultGinMiddlewares 默认的 gin server 使用的中间件列表
func DefaultGinMiddlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{
		// 记录请求处理日志，最顶层执行
		webserver.GinLogMiddleware(),
		// 捕获 panic 保存到 context 中由 GinLogger 统一打印， panic 时返回 500 JSON
		webserver.GinRecovery(response.Respond),
	}

	// 配置开启请求限频则添加限频中间件
	if viper.GetBool("ratelimiter.enable") {
		m = append(m, webserver.GinRatelimitMiddleware())
	}
	// i18n多语言
	if viper.GetBool("i18n.enable") {
		m = append(
			m,
			webserver.GinSetLanguage(),
		)
	}
	return m
}

func main() {
	configFile := flag.String("c", "./config.default.toml", "name of config file without format suffix)")
	flag.Parse()
	webserver.InitWithConfigFile(*configFile)

	// 依赖初始化
	models.Init()
	services.Init()

	// 创建 gin app
	middlewares := DefaultGinMiddlewares()
	app := webserver.NewGinEngine(middlewares...)
	// 注册路由
	routes.InitRouter(app)
	// 运行服务
	webserver.Run(app)
}
