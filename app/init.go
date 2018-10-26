// main init in proper order at here
package main

import (
	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/spf13/viper"
)

func init() {
	utils.InitViper([]utils.Option{
		utils.Option{"server.mode", "debug", "server mode: debug|test|release"},
		utils.Option{"server.bind", ":8080", "server bind address"},
		utils.Option{"log.level", "info", "log level: debug|info|warning|error|fatal|panic"},
		utils.Option{"log.formatter", "text", "log formatter: text|json"},
		utils.Option{"database.engine", "sqlite3", "database engine: mysql|postgres|sqlite3|mssql"},
		utils.Option{"database.address", "", "database address: host:port"},
		utils.Option{"database.name", "/tmp/gin-skeleton.db", "database name"},
		utils.Option{"database.username", "", "database username"},
		utils.Option{"database.password", "", "database password"},
	})

	utils.InitLogrus(viper.GetString("log.level"), viper.GetString("log.formatter"))
	utils.InitGormDB(viper.GetString("database.engine"), viper.GetString("database.address"), viper.GetString("database.name"), viper.GetString("database.username"), viper.GetString("database.password"))
}
