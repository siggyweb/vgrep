package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

// View creates the TUI representation
//
//goland:noinspection GoMixedReceiverTypes
func (m ShellModel) View() string {
	resultView := textStyle.Render(fmt.Sprintf("Result: %s", m.output))

	errorView := fmt.Sprintf("Error: %s", func() string {
		if m.err != nil {
			return m.err.Error()
		}
		return ""
	}())
	errorView = panelStyle.Render(errorView)

	inputView := panelStyle.Render(m.inputBuffer.View())

	return lipgloss.JoinVertical(lipgloss.Center, resultView, errorView, inputView)
}

// general panelStyle
var panelStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#0000FF")).
	Width(80).
	MaxWidth(75).
	//Height().
	// Padding(1).
	Border(lipgloss.NormalBorder(), true)

var textStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#6B4EFF")).       // Purple background
	Foreground(lipgloss.Color("#FFFFFF")).       // White text
	Border(lipgloss.InnerHalfBlockBorder()).     // Optional border
	BorderForeground(lipgloss.Color("#3333FF")). // Darker blue border
	Align(lipgloss.Left).                        // Center-align text
	Width(80).                                   // Set fixed width for wrapping
	MaxWidth(50)

	// todo function for error message formatting
