package main

import (
	"flag"
	"os"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/apis"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-gonic/gin"
)

// middlewares 返回 server 中需要使用的 gin 中间件
func middlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{}
	return m
}

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
	webserver.InitViperConfig(*configPath, *configName, *configType)

	// 创建 gin app
	app := webserver.NewGinEngine(middlewares()...)
	// 运行服务
	webserver.Run(app, apis.Register)
}
