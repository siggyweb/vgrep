package test

import (
	"github.com/siggyweb/vgrep/internal/stats"
	"github.com/siggyweb/vgrep/internal/tui"
	log "github.com/sirupsen/logrus/hooks/test"
	"testing"
)

// tui model test?

//func TestUpdateLoopOnKeypress(t *testing.T) {
//	// simulate user typing pwd command
//	testModel := CreateTestModel()
//	keysPressed := []rune{'p', 'w', 'd'}
//
//	var capturedModel tea.Model
//	var cmd tea.Cmd
//
//	for _, key := range keysPressed {
//		msg := tea.KeyMsg{
//			Type:  tea.KeyRunes,
//			Runes: []rune{key},
//			Alt:   false,
//			Paste: false,
//		}
//		capturedModel, cmd = testModel.Update(msg)
//		time.Sleep(600 * time.Millisecond)
//	}
//
//	time.Sleep(1000 * time.Millisecond)
//	testView := capturedModel.View()
//	expectedView := " Result: test/dir\nError:\ntest/dir>>"
//	if testView != expectedView {
//		t.Errorf("got view: \n %s \n expected view \n %s", testView, expectedView)
//	}
//}

func TestCommandValidator(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{
			"ls",
			true,
		},
		{
			input: "rm",
			valid: false,
		},
		{
			input: "grep",
			valid: true,
		},
		{
			input: "chmod",
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			if output := tui.ValidateCommand(test.input); output != test.valid {
				t.Errorf("want %v for binary %s, got %v", test.valid, test.input, output)
			}
		})
	}
}

func CreateTestModel() tui.ShellModel {
	testLogger, _ := log.NewNullLogger()

	var testModel = tui.ShellModel{
		CurrentDirectory: "test/dir>>",
		DebounceTag:      0,
		Err:              nil,
		Height:           100,
		History:          &tui.History{},
		InputBuffer:      tui.CreateInputBuffer("test/dir>>"),
		Logger:           testLogger,
		Output:           "",
		Stats:            &stats.SessionStatsModel{},
	}
	testModel.InputBuffer.Prompt = testModel.CurrentDirectory

	return testModel
}
