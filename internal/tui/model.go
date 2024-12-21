package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
	"os"
	"path/filepath"
	"strings"
)

// ShellModel represents the dynamic layer above the terminal which handles the interaction with the system shell below
// it implements the bubble tea application state model for terminal
type ShellModel struct {
	currentDirectory string
	err              error
	inputBuffer      textinput.Model
	logger           *log.Logger
	output           string // do I need a builder here?
	height           int
	width            int
	debounceTag      int
}

// InitialModel creates the starting state for the event loop
func InitialModel(logger *log.Logger) ShellModel {
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
		logger:           logger,
	}
	logger.Debugln("TUI state initialised")

	return model
}

// Init kicks off the event loop
func (m ShellModel) Init() tea.Cmd {
	return m.inputBuffer.Focus()
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
