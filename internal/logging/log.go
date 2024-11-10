package logging

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// ConfigureLogging sets up an instance of the logrus logger to be passed into the main model for writing logs to file.
// Logs are based around tea.Msg handling as these are the currency of the system and drive all behaviour.
// TickMsg are ignored by logging as their existence is inferred by the presence of CommandResponseMessage which cannot
// trigger without a tick
func ConfigureLogging() (*log.Logger, func() error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("could not retrieve home directory")
	}

	logPath := filepath.Join(homeDir, "app.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Could not open log file.")
	}

	logger := log.New()
	logger.SetOutput(logFile)
	logger.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

	return logger, logFile.Close
}
