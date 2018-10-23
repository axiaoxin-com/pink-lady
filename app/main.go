package main

import (
	"github.com/fvbock/endless"
	"github.com/spf13/viper"
)

func main() {
	app := initAPP()
	endless.ListenAndServe(viper.GetString("bind"), app)
}
