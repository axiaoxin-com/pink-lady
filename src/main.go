//go:generate swag init --dir ./ --generalInfo apis/apis.go --propertyStrategy snakecase --output ./apis/docs

package main

import (
	"flag"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/apis"
	"github.com/axiaoxin-com/pink-lady/services"
	"github.com/axiaoxin-com/pink-lady/webserver"
)

func main() {
	configFile := flag.String("c", "./config.default.toml", "name of config file without format suffix)")
	flag.Parse()
	webserver.InitWithConfigFile(*configFile)

	// 初始化或加载外部依赖服务
	if err := services.Init(); err != nil {
		logging.Error(nil, "services init error:"+err.Error())
	}

	// 创建 gin app
	middlewares := webserver.DefaultGinMiddlewares()
	app := webserver.NewGinEngine(middlewares...)
	// 运行服务
	webserver.Run(app, apis.Register)
}
