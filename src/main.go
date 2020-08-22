package main

import (
	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/pink-lady/apis"
	"github.com/spf13/viper"
)

func main() {
	goutils.InitWebAppViperConfig()
	app := goutils.NewGinEngine(viper.GetString("server.mode"), viper.GetBool("server.pprof"))
	goutils.RunWebApp(app, apis.Register)

}
