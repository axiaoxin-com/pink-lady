package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/axiaoxin/pink-lady/app/apis"
	"github.com/axiaoxin/pink-lady/app/db"
	"github.com/axiaoxin/pink-lady/app/router"
	"github.com/axiaoxin/pink-lady/app/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	workdir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	if err := utils.InitViper(workdir, "config", "GIN",
		utils.ViperOption{Name: "server.mode", Default: "debug", Desc: "server mode: debug|test|release"},
		utils.ViperOption{Name: "server.bind", Default: ":4869", Desc: "server bind address"},
		utils.ViperOption{Name: "log.level", Default: "info", Desc: "log level: debug|info|warning|error|fatal|panic"},
		utils.ViperOption{Name: "log.formatter", Default: "text", Desc: "log formatter: text|json"},
		utils.ViperOption{Name: "database.engine", Default: "sqlite3", Desc: "database engine: mysql|postgres|sqlite3|mssql"},
		utils.ViperOption{Name: "database.address", Default: "", Desc: "database address: host:port"},
		utils.ViperOption{Name: "database.name", Default: "/tmp/pink-lady.db", Desc: "database name"},
		utils.ViperOption{Name: "database.username", Default: "", Desc: "database username"},
		utils.ViperOption{Name: "database.password", Default: "", Desc: "database password"},
		utils.ViperOption{Name: "database.maxIdleConns", Default: 2, Desc: "sets the maximum number of connections in the idle connection pool."},
		utils.ViperOption{Name: "database.maxOpenConns", Default: 0, Desc: "sets the maximum number of open connections to the database."},
		utils.ViperOption{Name: "database.connMaxLifeMinutes", Default: 0, Desc: "sets the maximum amount of time(minutes) a connection may be reused."},
		utils.ViperOption{Name: "database.logMode", Default: true, Desc: "show detailed sql log"},
		utils.ViperOption{Name: "database.autoMigrate", Default: true, Desc: "auto migrate database when server starting"},
		utils.ViperOption{Name: "redis.mode", Default: "single-instance", Desc: "redis mode: single-instance|sentinel|cluster"},
		utils.ViperOption{Name: "redis.address", Default: "localhost:6379", Desc: "redis address, multiple sentinel/cluster addresses are separated by commas"},
		utils.ViperOption{Name: "redis.password", Default: "", Desc: "redis password"},
		utils.ViperOption{Name: "redis.db", Default: 0, Desc: "redis default db"},
		utils.ViperOption{Name: "redis.master", Default: "", Desc: "redis sentinel master name"},
		utils.ViperOption{Name: "sentry.dsn", Default: "", Desc: "sentry dsn"},
		utils.ViperOption{Name: "sentry.onlyCrashes", Default: "", Desc: "sentry only send crash reporting"},
	); err != nil {
		logrus.Error(err)
	}

	utils.InitLogger(os.Stdout, viper.GetString("log.level"), viper.GetString("log.formatter"))
	if err := db.InitGorm(); err != nil {
		logrus.Error(err)
	}
	if err := utils.InitRedis(viper.GetString("redis.mode"), viper.GetString("redis.address"), viper.GetString("redis.password"), viper.GetInt("redis.db"), viper.GetString("redis.master")); err != nil {
		logrus.Error(err)
	}
}

func main() {
	logrus.Info("===================================== pink-lady =====================================")
	defer db.DBInstanceMap.Close()
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
