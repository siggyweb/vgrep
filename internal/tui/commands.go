package tui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"strings"
	"time"
)

// ValidateCommand provides a whitelist of commands which are safe to be run within the event loop
func ValidateCommand(executable string) bool {
	safeCommands := []string{"pwd", "ls", "grep", "find", "locate", "which", "awk"}
	for _, cmd := range safeCommands {
		if executable == cmd {
			return true
		}
	}
	return false
}

// CreateCommand constructs terminal commands from users input with safety checks
func (m *ShellModel) CreateCommand() (*exec.Cmd, context.CancelFunc) {
	arguments := strings.Fields(m.InputBuffer.Value())
	l := len(arguments)
	if l == 0 {
		return nil, nil
	}

	var command *exec.Cmd
	valid := ValidateCommand(arguments[0])
	if !valid {
		m.Stats.IncrementInvalidCommands()
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	switch l {
	case 1:
		command = exec.CommandContext(ctx, arguments[0])
	default:
		command = exec.CommandContext(ctx, arguments[0], arguments[1:]...)
	}

	return command, cancel
}

// RunCommand executes shell commands as a tea.Cmd and routes the results back into the event loop
//
//goland:noinspection GoMixedReceiverTypes
func (m *ShellModel) RunCommand() tea.Cmd {
	return func() tea.Msg {
		command, cancel := m.CreateCommand()
		if command == nil {
			return nil
		}
		defer cancel()

		m.History.AddCommand(command.String())
		output, err := command.Output()
		if err != nil {
			m.Stats.IncrementErrors()
			return CommandResponseMessage{
				result: "",
				err:    err,
			}
		} else {
			m.Stats.IncrementCommandsRun()
			return CommandResponseMessage{
				result: string(output),
				err:    nil,
			}
		}
	}
}
