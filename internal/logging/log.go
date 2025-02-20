package logging

import (
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"reflect"
	"slices"
)

// ConfigureLogging sets up an instance of the logrus logger for dependency injection.
// Logs are based around tea.Msg handling as these are the currency of the system and drive all behaviour.
// TickMsg are ignored by logging.
func ConfigureLogging() InternalLogger {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("could not retrieve home directory")
	}

	logPath := filepath.Join(homeDir, "app.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Could not open log file.")
	}

	baseLogger := log.New()
	baseLogger.SetOutput(logFile)
	baseLogger.SetLevel(log.DebugLevel)
	baseLogger.SetReportCaller(true)

	appLogger := &MessageLogger{
		Logger:     baseLogger,
		LogFile:    logFile,
		FilterList: []string{"BlinkMsg"},
	}

	return appLogger
}

type InternalLogger interface {
	CleanUp()
	LogMessage(message tea.Msg, level log.Level)
	Infof(format string, args ...interface{})
}

type MessageLogger struct {
	Logger     *log.Logger
	LogFile    *os.File
	FilterList []string
}

func (l *MessageLogger) LogMessage(message tea.Msg, level log.Level) {
	messageType := reflect.TypeOf(message).Name()
	if slices.Contains(l.FilterList, messageType) {
		return
	}
	l.Logger.WithField("message_type", messageType).Logf(level, "contents: %v", message)
}

func (l MessageLogger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *MessageLogger) CleanUp() {
	_ = l.LogFile.Close()
}
