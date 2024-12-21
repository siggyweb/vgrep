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

const DebounceDuration = time.Second

// Update handles core routing for messages flowing through the MVU pipeline
//
//goland:noinspection GoMixedReceiverTypes
func (m ShellModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	m.logger.Debugf("handling message of type: %s,  message: %+v", reflect.TypeOf(message), message)
	var messageCommand tea.Cmd

	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.HandleWindowSizeMsg(msg)
	case tea.KeyMsg:
		messageCommand = m.HandleKeyMessage(msg)
	case TickMsg:
		if int(msg) == m.debounceTag {
			messageCommand = m.CommandRunner()
		}
	case CommandResponseMessage:
		m.HandleCommandResponseMessage(msg)
	}

	// Manage the state of the ti bubble via its own mvu event loop
	var inputBufferCmd tea.Cmd
	m.inputBuffer, inputBufferCmd = m.inputBuffer.Update(message)
	if len(m.inputBuffer.Value()) == 0 {
		m.output = ""
		m.err = nil
	}

	return m, tea.Batch(messageCommand, inputBufferCmd)
}

func (m *ShellModel) HandleKeyMessage(msg tea.KeyMsg) tea.Cmd {
	m.debounceTag++
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit

	case "ctrl+q":
		clipboard.Write(clipboard.FmtText, []byte(m.output))
	}
	return tea.Tick(DebounceDuration, func(_ time.Time) tea.Msg {
		return TickMsg(m.debounceTag)
	})
}

func (m *ShellModel) HandleCommandResponseMessage(msg CommandResponseMessage) {
	if msg.err == nil {
		m.output = msg.result
		m.err = nil
	} else {
		m.output = ""
		m.err = msg.err
	}
}

func (m *ShellModel) HandleWindowSizeMsg(msg tea.WindowSizeMsg) {
	m.height = msg.Height
	m.width = msg.Width
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
func (m *ShellModel) CommandRunner() tea.Cmd {
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
