package tui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"strings"
	"time"
)

// validateCommand provides a whitelist of commands which are safe to be run within the event loop
func validateCommand(executable string) bool {
	safeCommands := []string{"pwd", "ls", "grep", "find", "locate", "which", "awk"}
	for _, cmd := range safeCommands {
		if executable == cmd {
			return true
		}
	}
	return false
}

// CommandCreator constructs terminal commands from users input with safety checks
func (m *ShellModel) CommandCreator() (*exec.Cmd, context.CancelFunc) {
	arguments := strings.Fields(m.inputBuffer.Value())
	l := len(arguments)
	if l == 0 {
		return nil, nil
	}

	var command *exec.Cmd
	valid := validateCommand(arguments[0])
	if !valid {
		m.stats.IncrementInvalidCommands()
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

// CommandRunner executes shell commands as a tea.Cmd and routes the results back into the event loop
//
//goland:noinspection GoMixedReceiverTypes
func (m *ShellModel) CommandRunner() tea.Cmd {
	return func() tea.Msg {
		command, cancel := m.CommandCreator()
		if command == nil {
			return nil
		}
		defer cancel()

		m.history.AddCommand(command.String())
		output, err := command.Output()
		if err != nil {
			m.stats.IncrementErrors()
			return CommandResponseMessage{
				result: "",
				err:    err,
			}
		} else {
			m.stats.IncrementCommandsRun()
			return CommandResponseMessage{
				result: string(output),
				err:    nil,
			}
		}
	}
}
