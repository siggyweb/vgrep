package vgrep

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// bubbletea application state Model
type Model struct {
	result      string
	inputBuffer textinput.Model
	err         error
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "begin searching..."
	ti.Prompt = ">>"
	ti.Focus()

	model := Model{
		result:      "",
		inputBuffer: ti,
		err:         nil,
	}
	return model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.inputBuffer, cmd = m.inputBuffer.Update(message)
	return m, cmd
}

func (m Model) View() string {
	view := fmt.Sprintf("Result: %s \n", func() string {
		if m.err == nil {
			return m.result
		}
		return ""
	}())
	view += m.inputBuffer.View()
	return view
}
