package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	Log.SetOutput(os.Stdout) // Ensure logs go to stdout
}
