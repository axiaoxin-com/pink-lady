package main

import (
	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/spf13/viper"
)

func init() {
	utils.InitViper([]utils.Option{
		utils.Option{"server.mode", "debug", "server mode: debug|test|release"},
		utils.Option{"server.bind", ":8080", "server bind address"},
		utils.Option{"log.out", "stdout", "log output: stdout|stderr"},
		utils.Option{"log.level", "info", "log level: debug|info|warning|error|fatal|panic"},
		utils.Option{"log.formatter", "text", "log formatter: text|json"},
	})

	utils.InitLogrus(viper.GetString("log.level"), viper.GetString("log.formatter"), viper.GetString("log.out"))

}
