//go:generate swag init --dir ./ --generalInfo apis/apis.go --propertyStrategy snakecase --output ./apis/docs

package main

import (
	"flag"
	"fmt"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/apis"
	"github.com/axiaoxin-com/pink-lady/apis/response"
	"github.com/axiaoxin-com/pink-lady/services"
	"github.com/axiaoxin-com/pink-lady/statics"
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
	return m
}

func main() {
	configFile := flag.String("c", "./config.default.toml", "name of config file without format suffix)")
	flag.Parse()
	webserver.InitWithConfigFile(*configFile)

	db, err := goutils.GormMySQL("localhost")
	fmt.Println("=====>", db, err)
	// 初始化或加载外部依赖服务
	if err := services.Init(); err != nil {
		logging.Error(nil, "services init error:"+err.Error())
	}

	// 创建 gin app
	middlewares := DefaultGinMiddlewares()
	app := webserver.NewGinEngine(&statics.Files, middlewares...)
	// 注册路由
	apis.Register(app)
	// 运行服务
	webserver.Run(app)
}
