package tui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

// Update handles the changes of state for the model
func (m ShellModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	msgType := reflect.TypeOf(message)
	if msgType != reflect.TypeOf(TickMsg{}) {
		m.logger.Debugf("handling message, type: %s , message: %+v", msgType, message)
	}

	switch msg := message.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+q":
			clipboard.Write(clipboard.FmtText, []byte(m.output))
		}

	case TickMsg:
		return m, tea.Batch(
			m.CommandRunner(),
			tickEvery(),
		)

	case CommandResponseMessage:
		if msg.err == nil {
			m.output = msg.result
			m.err = nil
		} else {
			m.output = ""
			m.err = msg.err
		}
	}

	// finally manage the state of the ti bubble via its own mvu event loop
	var cmd tea.Cmd
	m.inputBuffer, cmd = m.inputBuffer.Update(message)
	if len(m.inputBuffer.Value()) == 0 {
		m.output = ""
		m.err = nil
	}

	return m, cmd
}

// CommandCreator constructs terminal commands from users input with safety checks
func (m ShellModel) CommandCreator() (*exec.Cmd, context.CancelFunc) {
	arguments := strings.Fields(m.inputBuffer.Value())
	l := len(arguments)
	if l == 0 {
		return nil, nil
	}

	var command *exec.Cmd
	valid := validateCommand(arguments[0])
	if !valid {
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

// ValidateCommand provides a whitelist of commands which are safe to be run within the event loop
func validateCommand(executable string) bool {
	safeCommands := []string{"pwd", "ls", "grep", "find", "locate", "which", "awk"}
	for _, cmd := range safeCommands {
		if executable == cmd {
			return true
		}
	}
	return false
}

// CommandRunner executes shell commands in a goroutine using tea Cmd capability and routes the results back into the event loop
func (m ShellModel) CommandRunner() tea.Cmd {
	return func() tea.Msg {
		command, cancel := m.CommandCreator()
		// if command is invalid abandon here as we cannot call cancel()
		if command == nil {
			return nil
		}
		// else set up the command with cancellation token and execute
		defer cancel()

		output, err := command.Output()
		if err != nil {
			return CommandResponseMessage{
				result: "",
				err:    err,
			}
		} else {
			return CommandResponseMessage{
				result: string(output),
				err:    nil,
			}
		}
	}
}

// tickEvery is the driver for the refresh rate of results in the view
func tickEvery() tea.Cmd {
	return tea.Every(time.Millisecond*500, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
