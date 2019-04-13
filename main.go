package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	VERSION    = "0.2.0"
	BUILD      = "0.2.0"
	log        = logrus.New()
	MainIsDone = make(chan bool)
)
var MyConfig Config

func init() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(logrus.DebugLevel)
}

func main() {

	// init stuff

	MyConfig.Load("config.json")

	// start bot
	StartBot(&MyConfig)

	// waiting for shutdown signal
	<-MainIsDone

}
