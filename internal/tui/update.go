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
		m.logger.Debugf("handling message, type: %s , message: %s", msgType, message)
	}

	switch msg := message.(type) {
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

	// manage the state of the ti bubble via its own mvu event loop
	var cmd tea.Cmd
	m.inputBuffer, cmd = m.inputBuffer.Update(message)
	if len(m.inputBuffer.Value()) == 0 {
		m.output = ""
		m.err = nil
	}

	return m, cmd
}

// CommandCreator forms shell commands to be executed async
func (m ShellModel) CommandCreator() (*exec.Cmd, context.CancelFunc) {
	// split the raw cmd text from the users input into args and form an executable command
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	arguments := strings.Fields(m.inputBuffer.Value())
	var command *exec.Cmd

	l := len(arguments)
	switch l {
	case 0:
		cancel()
		return nil, nil
	case 1:
		command = exec.CommandContext(ctx, arguments[0])
	default:
		command = exec.CommandContext(ctx, arguments[0], arguments[1:]...)
	}

	return command, cancel
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
