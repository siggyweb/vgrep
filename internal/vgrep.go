package vgrep

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
	"os/exec"
	"strings"
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
		}

	case TickMsg:
		return m, tea.Batch(
			m.GrepFetcher(),
			tickEvery(),
		)

	case GrepMessage:
		if msg.err == nil {
			m.output = msg.result
			m.err = nil
		} else {
			m.output = ""
			m.err = msg.err
		}
	}

	// the ti bubble has its own mvu loop, reset when user deletes all input
	m.inputBuffer, cmd = m.inputBuffer.Update(message)
	if len(m.inputBuffer.Value()) == 0 {
		m.output = ""
		m.err = nil
	}

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

func (m Model) GrepFetcher() tea.Cmd {
	return func() tea.Msg {
		// split the raw cmd text into args and handle
		var command *exec.Cmd
		arguments := strings.Fields(m.inputBuffer.Value())

		l := len(arguments)
		switch l {
		case 0:
			return nil
		case 1:
			command = exec.Command(arguments[0])
		default:
			command = exec.Command(arguments[0], arguments[1:]...)
		}

		output, err := command.Output()
		if err != nil {
			return GrepMessage{
				result: "",
				err:    err,
			}
		} else {
			return GrepMessage{
				result: string(output),
				err:    nil,
			}
		}
	}
}

type GrepMessage struct {
	result string // could be []string depending on result? requires testing
	err    error
}

type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Every(time.Millisecond*500, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
