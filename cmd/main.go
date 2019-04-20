package main

import (
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/rbtracker"
)

func main() {

	// init stuff
	rbtracker.MyConfig.Load("configs/config.json")

	// start bot
	logger.Logger.Info("Starting bot in go routine...")
	go rbtracker.StartBot()

	// waiting for shutdown signal
	<-rbtracker.KillBot

	defer logger.Logger.Info("Stopping bot go routine!")
}
