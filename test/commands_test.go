package test

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/siggyweb/vgrep/internal/stats"
	"github.com/siggyweb/vgrep/internal/tui"
	log "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
	"time"
)

func TestUpdateLoopOnKeypress(t *testing.T) {
	t.Skip("unfinished.")
	// simulate user typing pwd command
	testModel := CreateTestModel("")
	keysPressed := []rune{'p', 'w', 'd'}

	var capturedModel tea.Model
	for _, key := range keysPressed {
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{key},
			Alt:   false,
			Paste: false,
		}
		capturedModel, _ = testModel.Update(msg)
		fmt.Println(capturedModel.View())
		time.Sleep(600 * time.Millisecond)
	}

	testView := capturedModel.View()
	expectedView := " Result: test/dir\nError:\ntest/dir>>"
	if testView != expectedView {
		t.Errorf("got view: \n %s \n expected view \n %s", testView, expectedView)
	}
}

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

func TestCommandCreatorFormsCommands(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	tests := []struct {
		name      string
		inputText string
		expected  *exec.Cmd
	}{
		{
			"single-word command",
			"ls",
			exec.CommandContext(ctx, "ls"),
		},
		{
			"multi-word command",
			"grep Run",
			exec.CommandContext(ctx, "grep", "Run"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			model := CreateTestModel(test.inputText)
			actual, cancel := model.CommandCreator()
			defer cancel()

			if actual.Path != test.expected.Path {
				t.Errorf("got %v but expected %s", actual, test.expected)
			}
		})
	}
}

func TestCommandCreatorRejectsEmptyCommands(t *testing.T) {
	input := struct {
		name            string
		commandInput    string
		expectedCommand *exec.Cmd
	}{
		"An empty command text produces no shell command",
		"",
		nil,
	}
	t.Run(input.name, func(t *testing.T) {
		model := CreateTestModel(input.commandInput)
		actualCommand, cancel := model.CommandCreator()

		assert.Nil(t, actualCommand, "got %v but expected nil", actualCommand)
		assert.Nil(t, cancel, "got %v but expected nil", cancel)
	})
}

func CreateTestModel(input string) tui.ShellModel {
	testLogger, _ := log.NewNullLogger()

	textInput := textinput.Model{
		Prompt: "test/dir>>",
	}
	textInput.SetValue(input)

	var testModel = tui.ShellModel{
		CurrentDirectory: "test/dir",
		DebounceTag:      0,
		Err:              nil,
		Height:           100,
		Width:            100,
		History:          &tui.History{},
		InputBuffer:      textInput,
		Logger:           testLogger,
		Output:           "",
		Stats:            &stats.SessionStatsModel{},
	}

	return testModel
}
