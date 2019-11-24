package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axiaoxin/pink-lady/app/apis"
	"github.com/axiaoxin/pink-lady/app/db"
	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/axiaoxin/pink-lady/app/router"
	"github.com/axiaoxin/pink-lady/app/utils"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal("[FATAL] ", err)
	}
	if err := utils.InitViper(workdir, "config", "GIN",
		utils.ViperOption{Name: "server.mode", Default: "debug", Desc: "server mode: debug|test|release"},
		utils.ViperOption{Name: "server.bind", Default: ":4869", Desc: "server bind address"},
		utils.ViperOption{Name: "log.level", Default: "info", Desc: "log level: debug|info|warning|error|fatal|panic"},
		utils.ViperOption{Name: "log.formatter", Default: "text", Desc: "log formatter: text|json"},
	); err != nil {
		log.Println("[ERROR]", err)
	}

	if err := logging.InitLogger(); err != nil {
		log.Println("[ERROR] ", err)
	}
	if err := utils.InitSentry(); err != nil {
		log.Println("[ERROR] ", err)
	}
	if err := db.InitGorm(); err != nil {
		log.Println("[ERROR] ", err)
	}
	if err := utils.InitRedis(); err != nil {
		log.Println("[ERROR] ", err)
	}
}

func main() {
	log.Println("[INFO] ======================= pink-lady =======================")
	defer db.InstanceMap.Close()

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

	app := router.SetupRouter()
	apis.RegisterRoutes(app)
	bind := viper.GetString("server.bind")
	utils.EndlessServe(bind, app)
}
