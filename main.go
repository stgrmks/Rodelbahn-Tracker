package main

import (
	"github.com/stgrmks/Rodelbahn-Tracker/internal/config"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/slackbot"
)

var log = logger.Logger.WithField("package", "main")

var (
	VERSION = "0.2.0"
	BUILD   = "0.2.0"
	bot     slackbot.Bot
)

func init() {
	var MyConfig config.Config
	MyConfig.Load("configs/config.json")
	var crawlParams = slackbot.Param{
		Name:     "silent",
		Value:    "",
		Optional: true,
		Flag:     true,
	}
	var commands = []slackbot.Command{
		{
			Name:        "version",
			Description: "Show version and build number",
			ParamMap:    nil,
		},
		{
			Name:        "crawlNow",
			Description: "Start the Rodelbahn crawl",
			ParamMap:    map[string]slackbot.Param{crawlParams.Name: crawlParams},
		},
		{
			Name:        "showConfig",
			Description: "Show current config",
			ParamMap:    nil,
		},
		{
			Name:        "periodicCrawl",
			Description: "Start periodic crawls",
			ParamMap:    nil,
		},
	}
	var commandMap = make(map[string]slackbot.Command)
	for _, command := range commands {
		commandMap[command.Name] = command
	}
	bot = slackbot.Bot{
		Version:           VERSION,
		Build:             BUILD,
		Shutdown:          make(chan bool),
		StopPeriodicCrawl: make(chan bool),
		MyConfig:          MyConfig,
		CommandMap:        commandMap,
	}

	log.Info("Initated all structs and attributes")

}

func main() {

	defer log.Info("Stopping slackbot go routine!")

	log.Info("Starting slackbot in go routine...")
	go bot.StartBot()

	// waiting for shutdown signal
	<-bot.Shutdown

}
