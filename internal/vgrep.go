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
	return Model{
		result:      "",
		inputBuffer: textinput.New(),
		err:         nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
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
