package main

import (
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/rbtracker"
)

var log = logger.Logger.WithField("package", "main")

func main() {

	defer log.Info("Stopping bot go routine!")

	// init stuff
	rbtracker.MyConfig.Load("configs/config.json")

	// start bot
	log.Info("Starting bot in go routine...")
	go rbtracker.StartBot()

	// waiting for shutdown signal
	<-rbtracker.KillBot

}
