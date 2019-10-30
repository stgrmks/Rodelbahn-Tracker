package main

import (
	"github.com/stgrmks/Rodelbahn-Tracker/internal/bot"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
)

var log = logger.Logger.WithField("package", "main")

func main() {

	defer log.Info("Stopping bot go routine!")

	// init stuff
	bot.MyConfig.Load("configs/config.json")

	// start bot
	log.Info("Starting bot in go routine...")
	go bot.StartBot()

	// waiting for shutdown signal
	<-bot.KillBot

}
