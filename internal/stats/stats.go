package stats

import (
	"fmt"
	"time"
)

// StatCollector represents an object that collects basic metrics about the tea app during runtime and publishes
// to the logfile before the program exits
type StatCollector interface {
	Init()
	IncrementCommandsRun()
	IncrementInvalidCommands()
	IncrementErrors()
	GetSummary() string
}

type SessionStatsModel struct {
	startTime           time.Time
	commandCount        int
	errors              int
	invalidCommandCount int
}

func (m *SessionStatsModel) Init() {
	m.startTime = time.Now()
	m.commandCount = 0
	m.errors = 0
}

// IncrementInvalidCommands tracks the number of shell commands which did not pass validation
func (m *SessionStatsModel) IncrementInvalidCommands() {
	m.invalidCommandCount++
}

// IncrementCommandsRun tracks the number of shell commands which were executed without error
func (m *SessionStatsModel) IncrementCommandsRun() {
	m.commandCount++
}

// IncrementErrors tracks the number of shell commands which were executed but produced an error
func (m *SessionStatsModel) IncrementErrors() {
	m.errors++
}

// GetSummary produces a report output of the counter totals tracked during this session
func (m *SessionStatsModel) GetSummary() string {
	timeElapsed := time.Since(m.startTime).Seconds()
	result := fmt.Sprintf("Stats summary:  Duration:%fs InvalidCommands:%d CommandsExecuted:%d  Errors:%d",
		timeElapsed, m.invalidCommandCount, m.commandCount, m.errors)
	return result
}
