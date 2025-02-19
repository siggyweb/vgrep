package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/siggyweb/vgrep/internal/stats"
	log "github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
	"os"
	"path/filepath"
	"strings"
)

// ShellModel represents the dynamic layer above the terminal which handles the interaction with the system shell below
// it implements the bubble tea application state model for the user's terminal
type ShellModel struct {
	CurrentDirectory string
	DebounceTag      int
	Err              error
	Height           int
	Width            int
	History          HistoryModel
	InputBuffer      textinput.Model
	Logger           *log.Logger
	Output           string
	Stats            stats.StatCollector
}

// InitialModel creates the starting state for the event loop
func InitialModel(logger *log.Logger, statsModel stats.StatCollector) ShellModel {
	workingDirectory, err := FetchWorkingDirectory()
	if err != nil {
		fmt.Println("could not obtain current working directory, safely quitting.")
		tea.Quit()
	}
	workingDirectory = filepath.Base(workingDirectory)

	ti := CreateInputBuffer(workingDirectory)

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	statsModel.Init()

	shellHistory := &History{
		storedCommands: nil,
		index:          0,
	}

	model := ShellModel{
		CurrentDirectory: workingDirectory,
		Output:           "",
		InputBuffer:      ti,
		Err:              nil,
		Logger:           logger,
		Stats:            statsModel,
		History:          shellHistory,
	}
	logger.Debugln("TUI state initialised")

	return model
}

// Init kicks off the event loop
//
//goland:noinspection GoMixedReceiverTypes
func (m ShellModel) Init() tea.Cmd {
	return m.InputBuffer.Focus()
}

func CreateInputBuffer(workingDir string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "begin searching..."
	ti.Prompt = workingDir + ">>"
	ti.Focus()

	return ti
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
