package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// Init initializes the logger
func Init(level, logFile string) error {
	Log = logrus.New()

	// Set log level
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	Log.SetLevel(lvl)

	// Set formatter
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set output (both file and stdout)
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		Log.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		Log.SetOutput(os.Stdout)
	}

	return nil
}
