package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
	"reflect"
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

	// Manage the state of the textinput bubble via its own mvu event loop
	var inputBufferCmd tea.Cmd
	m.inputBuffer, inputBufferCmd = m.inputBuffer.Update(message)
	if len(m.inputBuffer.Value()) == 0 {
		m.output = ""
		m.err = nil
	}

	return m, tea.Batch(messageCommand, inputBufferCmd)
}

// HandleKeyMessage manages translation of input from keyboard to keyboard shortcuts, running shell commands
// and command history. Debounced using int tags in the ShellModel.
func (m *ShellModel) HandleKeyMessage(msg tea.KeyMsg) tea.Cmd {
	m.debounceTag++
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit
	case "ctrl+q":
		clipboard.Write(clipboard.FmtText, []byte(m.output))
	case "up":
		m.inputBuffer.SetValue(m.history.GetPreviousCommand())
	case "down":
		m.inputBuffer.SetValue(m.history.GetNextCommand())
	}

	return tea.Tick(DebounceDuration, func(_ time.Time) tea.Msg {
		return TickMsg(m.debounceTag)
	})
}

// HandleCommandResponseMessage receives the results of valid shell commands executed directly via binaries, into the view.
func (m *ShellModel) HandleCommandResponseMessage(msg CommandResponseMessage) {
	if msg.err == nil {
		m.output = msg.result
		m.err = nil
	} else {
		m.output = ""
		m.err = msg.err
	}
}

// HandleWindowSizeMsg resizes tracks the dimensions of the TUI after a resize and stores in the ShellModel.
func (m *ShellModel) HandleWindowSizeMsg(msg tea.WindowSizeMsg) {
	m.height = msg.Height
	m.width = msg.Width
}
