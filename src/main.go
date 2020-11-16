package main

import (
	"flag"
	"os"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/apis"
	"github.com/axiaoxin-com/pink-lady/services"
	"github.com/axiaoxin-com/pink-lady/webserver"
)

func main() {
	// 根据命令行参数加载配置文件到 viper
	workdir, err := os.Getwd()
	if err != nil {
		logging.Warn(nil, "get workdir failed:"+err.Error())
		workdir = "."
	}
	configPath := flag.String("p", workdir, "path of config file")
	configName := flag.String("c", "config.default", "name of config file without format suffix)")
	configType := flag.String("t", "toml", "type of config file format")
	flag.Parse()
	webserver.InitWithConfigFile(*configPath, *configName, *configType)

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
