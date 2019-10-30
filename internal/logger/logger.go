package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = logrus.New()

func init() {

	Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		ForceColors:   true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Logger.SetOutput(os.Stdout)

	// Only Logger the warning severity or above.
	Logger.SetLevel(logrus.InfoLevel)
}
