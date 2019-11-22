package utils

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ViperOption the InitViper option type
type ViperOption struct {
	Name    string
	Default interface{}
	Desc    string
}

// InitViper init viper by default value, ENV, cmd flag and config file
// you can use switch to reload server when config file changed
func InitViper(configPath, configName string, envPrefix string, options ...ViperOption) error {
	viper.SetEnvPrefix(envPrefix)
	for _, option := range options {
		// set default value
		viper.SetDefault(option.Name, option.Default)

		// bind ENV
		viper.BindEnv(option.Name)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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

	// load conf file
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "viper read in config error")
	}
	log.Printf("loaded %s in %s\n", configName, configPath)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.ReadInConfig()
		log.Println("Config file changed, read in config:", e.Name)
	})
	return nil
}
