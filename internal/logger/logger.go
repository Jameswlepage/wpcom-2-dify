package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Log is a global logger instance.
var Log = logrus.New()

func init() {
	// Set logging level and format
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	// Configure output to both stdout and a file.
	// The file will contain raw requests/responses.
	// Make sure the directory /app exists in your Dockerfile/build environment, or adjust as needed.
	logFilePath := "/app/requests.log"
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// If we cannot open the file, fallback to stdout only.
		Log.Warnf("Could not open log file %s for writing: %v", logFilePath, err)
		Log.SetOutput(os.Stdout)
		return
	}

	// Write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, file)
	Log.SetOutput(multiWriter)
}
