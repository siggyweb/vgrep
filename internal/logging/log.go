package logging

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
)

func ConfigureLogging() func() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("could not retrieve home directory")
	}

	logPath := filepath.Join(homeDir, "app.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Could not open log file.")
	}

	logger := logrus.StandardLogger()
	logger.SetOutput(logFile)

	return logFile.Close
}
