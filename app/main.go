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
