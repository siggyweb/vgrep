package test

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/siggyweb/vgrep/internal/stats"
	"github.com/siggyweb/vgrep/internal/tui"
	log "github.com/sirupsen/logrus/hooks/test"
	"testing"
)

// tui model test?

func TestCommandCreatorCreatesValidCommand(t *testing.T) {

}

func CreateTestModel() tui.ShellModel {
	testLogger, _ := log.NewNullLogger()

	var testModel = tui.ShellModel{
		CurrentDirectory: "test/dir",
		DebounceTag:      0,
		Err:              nil,
		Height:           100,
		History:          &tui.History{},
		InputBuffer:      textinput.Model{},
		Logger:           testLogger,
		Output:           "",
		Stats:            &stats.SessionStatsModel{},
	}

	return testModel
}
