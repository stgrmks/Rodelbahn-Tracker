package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = logrus.New()

func init() {

	// Log as JSON instead of the default ASCII formatter.
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Logger.SetOutput(os.Stdout)

	// Only Logger the warning severity or above.
	Logger.SetLevel(logrus.DebugLevel)
}
