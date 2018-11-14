// gin-skeleton: Typically gin-based web application's organizational structure
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/axiaoxin/gin-skeleton/app/apis"
	"github.com/axiaoxin/gin-skeleton/app/apis/router"
	"github.com/axiaoxin/gin-skeleton/app/models"
	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	utils.InitViper("app", "GIN", []utils.ViperOption{
		utils.ViperOption{"server.mode", "debug", "server mode: debug|test|release"},
		utils.ViperOption{"server.bind", ":9090", "server bind address"},
		utils.ViperOption{"log.level", "info", "log level: debug|info|warning|error|fatal|panic"},
		utils.ViperOption{"log.formatter", "text", "log formatter: text|json"},
		utils.ViperOption{"database.engine", "sqlite3", "database engine: mysql|postgres|sqlite3|mssql"},
		utils.ViperOption{"database.address", "", "database address: host:port"},
		utils.ViperOption{"database.name", "/tmp/gin-skeleton.db", "database name"},
		utils.ViperOption{"database.username", "", "database username"},
		utils.ViperOption{"database.password", "", "database password"},
		utils.ViperOption{"database.maxIdleConns", 2, "sets the maximum number of connections in the idle connection pool."},
		utils.ViperOption{"database.maxOpenConns", 0, "sets the maximum number of open connections to the database."},
		utils.ViperOption{"database.connMaxLifeMinutes", 0, "sets the maximum amount of time(minutes) a connection may be reused."},
		utils.ViperOption{"database.logMode", true, "show detailed sql log"},
		utils.ViperOption{"database.autoMigrate", true, "auto migrate database when server starting"},
		utils.ViperOption{"redis.mode", "single-instance", "redis mode: single-instance|sentinel|cluster"},
		utils.ViperOption{"redis.address", "localhost:6379", "redis address, multiple sentinel/cluster addresses are separated by commas"},
		utils.ViperOption{"redis.password", "", "redis password"},
		utils.ViperOption{"redis.db", 0, "redis default db"},
		utils.ViperOption{"redis.master", "", "redis sentinel master name"},
		utils.ViperOption{"sentry.dsn", "", "sentry dsn"},
		utils.ViperOption{"sentry.onlyCrashes", "", "sentry only send crash reporting"},
	})

	utils.InitLogrus(os.Stdout, viper.GetString("log.level"), viper.GetString("log.formatter"))
	utils.InitGormDB(viper.GetString("database.engine"), viper.GetString("database.address"), viper.GetString("database.name"), viper.GetString("database.username"), viper.GetString("database.password"), viper.GetInt("database.maxIdleConns"), viper.GetInt("database.maxOpenConns"), viper.GetInt("database.connMaxLifeMinutes"), viper.GetBool("database.logMode"))
	if viper.GetBool("database.autoMigrate") {
		models.Migrate()
	}
	utils.InitRedis(viper.GetString("redis.mode"), viper.GetString("redis.address"), viper.GetString("redis.password"), viper.GetInt("redis.db"), viper.GetString("redis.master"))
}

func main() {
	defer utils.DB.Close()
	// TODO: imp in cli
	version := pflag.Bool("version", false, "show version")
	check := pflag.Bool("check", false, "check everything need to be checked")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	if *version {
		fmt.Println(apis.VERSION)
		os.Exit(0)
	}
	if *check {
		fmt.Println("I'm fine :)")
		os.Exit(0)
	}

	mode := strings.ToLower(viper.GetString("server.mode"))
	sentryDSN := viper.GetString("sentry.dsn")
	sentryOnlyCrashes := viper.GetBool("sentry.onlycrashes")
	app := router.SetupRouter(mode, sentryDSN, sentryOnlyCrashes)
	apis.RegisterRoutes(app)
	bind := viper.GetString("server.bind")
	utils.EndlessServe(bind, app)
}
