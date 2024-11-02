package vgrep

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

// bubbletea application state Model
type Model struct {
	result      string // do I need a builder here?
	inputBuffer textinput.Model
	err         error
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "begin searching..."
	ti.Prompt = ">>"
	ti.Focus()

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

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

	switch msg := message.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+q":
			clipboard.Write(clipboard.FmtText, []byte(m.result))

			// for testing purposes to see key input received
			//default:
			//	m.result += msg.String()
		}
	}

	// the ti bubble has its own mvu loop
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

type GrepMessage struct {
	result string // could be []string depending on result?
	err    error
}
