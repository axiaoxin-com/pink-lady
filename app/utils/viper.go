package utils

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
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
func InitViper(configName string, envPrefix string, options []ViperOption) error {
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
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("/etc")
	err := viper.ReadInConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// TODO
		logrus.Debug("TODO: reload gin server when config changed by swicther")
	})
	return err
}
