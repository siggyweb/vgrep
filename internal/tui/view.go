package tui

import "fmt"

// View creates the TUI representation
//
//goland:noinspection GoMixedReceiverTypes
func (m ShellModel) View() string {
	view := fmt.Sprintf("Result: %s \n", m.Output)
	view += fmt.Sprintf("Error: %s \n", func() string {
		if m.Err != nil {
			return m.Err.Error()
		}
		return ""
	}())
	view += m.InputBuffer.View()
	return view
}
