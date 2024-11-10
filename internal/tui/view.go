package tui

import "fmt"

// View creates the TUI representation
func (m ShellModel) View() string {
	// todo wrap lines + apply styling
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
