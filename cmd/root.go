package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	Host             string   `mapstructure:"host"`
	Database         string   `mapstructure:"database"`
	User             string   `mapstructure:"user"`
	Password         string   `mapstructure:"password"`
	Collection       string   `mapstructure:"collection"`
	BaseURL          string   `mapstructure:"baseURL"`
	ExtURL           string   `mapstructure:"extURL"`
	RbList           []string `mapstructure:"RbList"`
	SlackWebHook     string   `mapstructure:"SlackWebHook"`
	SlackBotToken    string   `mapstructure:"SlackBotToken"`
	Notify           bool     `mapstructure:"Notify"`
	Cron             string   `mapstructure:"Cron"`
	ActiveSession    *mgo.Session
	ActiveCollection *mgo.Collection
}

var (
	cfgFile string
	config  Config
	RootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  "",

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// if --config is passed, attempt to parse the config file
			if cfgFile != "" {

				// get the filepath
				abs, err := filepath.Abs(cfgFile)
				if err != nil {
					fmt.Printf("Error reading filepath: ", err.Error())
				}

				// get the config name
				base := filepath.Base(abs)

				// get the path
				path := filepath.Dir(abs)

				viper.SetConfigName(strings.Split(base, ".")[0])
				viper.AddConfigPath(path)

				// Find and read the config file; Handle errors reading the config file
				if err := viper.ReadInConfig(); err != nil {
					fmt.Printf("Failed to read config file: ", err.Error())
					os.Exit(1)
				} else {
					fmt.Printf("Using Config: %s\n", viper.ConfigFileUsed())
					if err := viper.Unmarshal(&config); err != nil {
						fmt.Println("Could not load config.")
					}
					log.Printf("%+v\n", config)
				}
				viper.WatchConfig()
				viper.OnConfigChange(func(e fsnotify.Event) {
					fmt.Println("Config file changed:", e.Name)
					if err := viper.Unmarshal(&config); err != nil {
						fmt.Println("Could not load config.")
					}
					log.Println(config)
				})
			} else {
				fmt.Println("Config is required!")
				os.Exit(1)
			}

			// overwrite slack hooks if set in ENV. for debugging
			for _, envVarString := range []string{"RBT_SlackWebHook", "RBT_SlackBotToken"} {
				envVar := os.Getenv(envVarString)
				if envVar != "" {
					fieldName := strings.Split(envVarString, "_")[1]
					rv := reflect.ValueOf(&config).Elem()
					fv := rv.FieldByName(fieldName)
					fv.SetString(envVar)
					log.Printf("Setting %s to %s", fieldName, envVar)
				}
			}

		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
)

func Execute(v, b string) {
	Version = v
	Build = b

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	// persistent flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.json", "config file")
	RootCmd.PersistentFlags().BoolVar(&config.Notify, "n", true, "turn notifications on")

	// commands
	RootCmd.AddCommand(version)
	RootCmd.AddCommand(crawl)
	RootCmd.AddCommand(initCrawl)
	RootCmd.AddCommand(periodicCrawl)
	RootCmd.AddCommand(bot)

}
