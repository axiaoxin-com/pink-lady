package main

import (
	"fmt"
	"os"

	"github.com/axiaoxin/gin-skeleton/app/handlers"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Option struct {
	Name    string
	Default interface{}
	Desc    string
}

func initViper() {
	options := []Option{
		Option{"mode", "debug", "server mode: debug|test|release"},
		Option{"bind", ":8080", "server bind address"},
		Option{"log.out", "stdout", "log output: stdout|stderr"},
		Option{"log.level", "info", "log level: debug|info|warning|error|fatal|panic"},
		Option{"log.formatter", "text", "log formatter: text|json"},
	}

	viper.SetEnvPrefix("GIN")
	for _, option := range options {
		// set default value
		viper.SetDefault(option.Name, option.Default)

		// bind ENV
		viper.BindEnv(option.Name)

		// cmd
		switch option.Default.(type) {
		case int:
			pflag.Int(option.Name, option.Default.(int), option.Desc)
		case string:
			pflag.String(option.Name, option.Default.(string), option.Desc)
		case bool:
			pflag.Bool(option.Name, option.Default.(bool), option.Desc)

		}
	}
	// TODO: imp in cli
	version := pflag.Bool("version", false, "show version")
	check := pflag.Bool("check", false, "check everything need to be checked")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	if *version {
		fmt.Println(handlers.VERSION)
		os.Exit(0)
	}
	if *check {
		fmt.Println("I'm fine :)")
		os.Exit(0)
	}

	// load conf file
	viper.SetConfigName("app")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		initAPP()
	})
}
