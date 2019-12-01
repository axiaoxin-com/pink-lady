package main

import (
	"log"
	"os"

	"pink-lady/app/apis"
	"pink-lady/app/router"
	"pink-lady/app/utils"

	"github.com/spf13/viper"
)

func main() {
	log.Println("[INFO] ============ pink-lady ============")
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal("[FATAL] ", err)
	}
	app := router.SetupRouter(workdir, "config")
	apis.RegisterRoutes(app)
	bind := viper.GetString("server.bind")
	utils.EndlessServe(bind, app)
}
