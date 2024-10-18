package vgrep

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

// bubbletea application state Model
type Model struct {
	resultdisplay string
	inputBuffer   string
	err           error
}

func InitialModel() Model {
	return Model{
		// our display is a text box
		resultdisplay: "",
		inputBuffer:   "",
		err:           nil,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
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
	view := ""
	view += fmt.Sprintf("Result: %s \n", func() string {
		if m.err == nil {
			return m.resultdisplay
		}
		return ""
	}())
	view += fmt.Sprintf("Current Command Text: %v \n", m.inputBuffer)
	return view
}
