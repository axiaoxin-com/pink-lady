package utils

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ViperOption the InitViper option type
type ViperOption struct {
	Name    string
	Default interface{}
	Desc    string
}

// NewViperOption return ViperOption pointer
func NewViperOption(name string, defaultValue interface{}, desc string) ViperOption {
	return ViperOption{
		Name:    name,
		Default: defaultValue,
		Desc:    desc,
	}
}

// InitViper init viper by default value, ENV, cmd flag and config file
// you can use switch to reload server when config file changed
func InitViper(configpaths []string, configname string, envPrefix string, options ...ViperOption) error {
	viper.SetEnvPrefix(envPrefix)
	for _, option := range options {
		// set default value
		viper.SetDefault(option.Name, option.Default)

		// bind ENV
		viper.BindEnv(option.Name)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	// load conf file
	viper.SetConfigName(configname)
	for _, configpath := range configpaths {
		viper.AddConfigPath(configpath)
	}
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "viper read in config error")
	}
	log.Printf("[INFO] loaded %s in %v\n", configname, configpaths)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.ReadInConfig()
		log.Println("[WARN] Config file changed, read in config:", e.Name)
	})
	return nil
}
