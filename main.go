package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("Welcome vgrep")
}

// bubbletea application state model
type model struct {
	display     string
	inputBuffer string
	err         error
}

func intialModel() model {
	return model{
		// our display is a text box
		display:     "",
		inputBuffer: "",
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return ""
}
