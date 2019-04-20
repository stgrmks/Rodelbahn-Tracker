package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	Host          string   `mapstructure:"host"`
	Database      string   `mapstructure:"database"`
	User          string   `mapstructure:"user"`
	Password      string   `mapstructure:"password"`
	Collection    string   `mapstructure:"collection"`
	BaseURL       string   `mapstructure:"baseURL"`
	ExtURL        string   `mapstructure:"extURL"`
	RbList        []string `mapstructure:"RbList"`
	SlackWebHook  string   `mapstructure:"SlackWebHook"`
	SlackBotToken string   `mapstructure:"SlackBotToken"`
	Notify        bool     `mapstructure:"Notify"`
	Cron          string   `mapstructure:"Cron"`
}

func (c *Config) Load(cfgPath string) {
	if cfgPath != "" {

		// get the filepath
		abs, err := filepath.Abs(cfgPath)
		if err != nil {
			fmt.Printf("Error reading filepath: ", err.Error())
		}

		// get the c name
		base := filepath.Base(abs)

		// get the path
		path := filepath.Dir(abs)

		viper.SetConfigName(strings.Split(base, ".")[0])
		viper.AddConfigPath(path)

		// Find and read the c file; Handle errors reading the c file
		if err := viper.ReadInConfig(); err != nil {
			logger.Logger.Infof("Failed to read c file: ", err.Error())
			os.Exit(1)
		} else {
			logger.Logger.Infof("Using Config: %s", viper.ConfigFileUsed())
			if err := viper.Unmarshal(c); err != nil {
				logger.Logger.Info("Could not load config!")
			}
			logger.Logger.Debugf("Config: %s", c)
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			logger.Logger.Infof("Config file changed:", e.Name)
			if err := viper.Unmarshal(c); err != nil {
				logger.Logger.Warn("Could not load config!")
			}
		})
	} else {
		logger.Logger.Fatal("Config is required!")
	}
	// overwrite slack hooks if set in ENV. for debugging
	for _, envVarString := range []string{"RBT_SlackWebHook", "RBT_SlackBotToken"} {
		envVar := os.Getenv(envVarString)
		if envVar != "" {
			fieldName := strings.Split(envVarString, "_")[1]
			rv := reflect.ValueOf(c).Elem()
			fv := rv.FieldByName(fieldName)
			fv.SetString(envVar)
			logger.Logger.Infof("Setting %s to %s", fieldName, envVar)
		}
	}

}
