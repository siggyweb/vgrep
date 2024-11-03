package vgrep

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
	"os/exec"
	"time"
)

// bubbletea application state model
type Model struct {
	output      string // do I need a builder here?
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
		output:      "",
		inputBuffer: ti,
		err:         nil,
	}
	return model
}

func (m Model) Init() tea.Cmd {
	return tickEvery()
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := message.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+q":
			clipboard.Write(clipboard.FmtText, []byte(m.output))

			// for testing purposes to see key input received
			//default:
			//	m.output += msg.String()
		}

	case TickMsg:
		//m.output += "tick"
		return m, tea.Batch(
			m.GrepFetcher(),
			tickEvery(),
		)

	case GrepMessage:
		if msg.err == nil {
			m.output = msg.result
		} else {
			m.output = msg.err.Error()
		}
	}

	// the ti bubble has its own mvu loop
	m.inputBuffer, cmd = m.inputBuffer.Update(message)
	return m, cmd
}

func (m Model) View() string {
	view := fmt.Sprintf("Result: %s \n", m.output)
	view += fmt.Sprintf("Error: %s \n", func() string {
		if m.err != nil {
			return m.err.Error()
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

func (m Model) GrepFetcher() tea.Cmd {
	return func() tea.Msg {
		command := exec.Command(m.inputBuffer.View())
		output, err := command.Output()
		if err == nil {
			return GrepMessage{
				result: string(output),
				err:    nil,
			}
		}
		return GrepMessage{
			result: "",
			err:    err,
		}
	}
}

type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
