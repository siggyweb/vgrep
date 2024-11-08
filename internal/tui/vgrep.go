package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ShellModel represents the dynamic layer above the terminal which handles the interaction with the system shell below
// it implements the bubble tea application state model for terminal
type ShellModel struct {
	currentDirectory string
	err              error
	inputBuffer      textinput.Model
	output           string // do I need a builder here?
}

// InitialModel creates the starting state for the event loop
func InitialModel() ShellModel {
	workingDirectory, err := FetchWorkingDirectory()
	if err != nil {
		fmt.Println("could not obtain current working directory, quitting")
		tea.Quit()
	}
	workingDirectory = filepath.Base(workingDirectory)

	ti := textinput.New()
	ti.Placeholder = "begin searching..."
	ti.Prompt = workingDirectory + ">>"
	ti.Focus()

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	model := ShellModel{
		currentDirectory: workingDirectory,
		output:           "",
		inputBuffer:      ti,
		err:              nil,
	}
	return model
}

// Init kicks off the event loop
func (m ShellModel) Init() tea.Cmd {
	return tea.Batch(tickEvery(), m.inputBuffer.Focus())
}

// Update handles the changes of state for the model
func (m ShellModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := message.(type) {
	case tea.KeyMsg:
		log.Debug(msg)
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+q":
			clipboard.Write(clipboard.FmtText, []byte(m.output))
		}

	case TickMsg:
		return m, tea.Batch(
			m.CommandRunner(),
			tickEvery(),
		)

	case CommandResponseMessage:
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

// View creates the TUI representation
func (m ShellModel) View() string {
	// todo split input based on \n and wrap lines.
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

// CommandCreator forms shell commands to be executed async
func (m ShellModel) CommandCreator() (*exec.Cmd, context.CancelFunc) {
	// split the raw cmd text from the users input into args and form an executable command
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	arguments := strings.Fields(m.inputBuffer.Value())
	var command *exec.Cmd

	l := len(arguments)
	switch l {
	case 0:
		cancel()
		return nil, nil
	case 1:
		command = exec.CommandContext(ctx, arguments[0])
	default:
		command = exec.CommandContext(ctx, arguments[0], arguments[1:]...)
	}

	return command, cancel
}

// CommandRunner executes shell commands in a goroutine using tea Cmd capability and routes the results back into the event loop
func (m ShellModel) CommandRunner() tea.Cmd {
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

// FetchWorkingDirectory Retrieves and formats the full path to the current working directory
func FetchWorkingDirectory() (string, error) {
	output, err := os.Getwd()
	if err != nil {
		return "", err
	}
	result := strings.TrimSpace(output)
	return result, nil
}

// tickEvery is the driver for the refresh rate of results in the view
func tickEvery() tea.Cmd {
	return tea.Every(time.Millisecond*500, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
