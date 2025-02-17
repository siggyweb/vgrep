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
	m.Logger.Debugf("handling message of type: %s,  message: %+v", reflect.TypeOf(message), message)
	var messageCommand tea.Cmd

	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.HandleWindowSizeMsg(msg)
	case tea.KeyMsg:
		messageCommand = m.HandleKeyMessage(msg)
	case TickMsg:
		if int(msg) == m.DebounceTag {
			messageCommand = m.CommandRunner()
		}
	case CommandResponseMessage:
		m.HandleCommandResponseMessage(msg)
	}

	// Manage the state of the textinput bubble via its own mvu event loop
	var inputBufferCmd tea.Cmd
	m.InputBuffer, inputBufferCmd = m.InputBuffer.Update(message)
	if len(m.InputBuffer.Value()) == 0 {
		m.Output = ""
		m.Err = nil
	}

	return m, tea.Batch(messageCommand, inputBufferCmd)
}

// HandleKeyMessage manages translation of input from keyboard to keyboard shortcuts, running shell commands
// and command history. Debounced using int tags in the ShellModel.
func (m *ShellModel) HandleKeyMessage(msg tea.KeyMsg) tea.Cmd {
	m.DebounceTag++
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit
	case "ctrl+q":
		clipboard.Write(clipboard.FmtText, []byte(m.Output))
	case "up":
		m.InputBuffer.SetValue(m.History.GetPreviousCommand())
	case "down":
		m.InputBuffer.SetValue(m.History.GetNextCommand())
	}

	return tea.Tick(DebounceDuration, func(_ time.Time) tea.Msg {
		return TickMsg(m.DebounceTag)
	})
}

// HandleCommandResponseMessage receives the results of valid shell commands executed directly via binaries, into the view.
func (m *ShellModel) HandleCommandResponseMessage(msg CommandResponseMessage) {
	if msg.err == nil {
		m.Output = msg.result
		m.Err = nil
	} else {
		m.Output = ""
		m.Err = msg.err
	}
}

// HandleWindowSizeMsg resizes tracks the dimensions of the TUI after a resize and stores in the ShellModel.
func (m *ShellModel) HandleWindowSizeMsg(msg tea.WindowSizeMsg) {
	m.Height = msg.Height
	m.Width = msg.Width
}
