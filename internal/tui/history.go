package tui

import "strings"

// HistoryModel is an interface where objects that implement the contract can provide backtracking capability to shell users.
type HistoryModel interface {
	GetPreviousCommand() string
	GetNextCommand() string
	AddCommand(string)
}

// History stores previously run shell commands which can be easily accessed, allowing users to backtrack effectively in their shell session.
type History struct {
	storedCommands []string
	index          int
}

// GetPreviousCommand fetches the next stored command from the shell model. Fetching from beyond either end of the list
// will return an empty string.
func (h *History) GetPreviousCommand() string {
	if len(h.storedCommands) == 0 {
		return ""
	} else if h.index > 0 {
		h.index--
		return h.storedCommands[h.index]
	}

	return ""
}

// GetNextCommand fetches the next stored command from the shell model. Fetching from beyond either end of the list
// will return an empty string.
func (h *History) GetNextCommand() string {
	if len(h.storedCommands) == 0 {
		return ""
	}

	if h.index < len(h.storedCommands)-1 {
		h.index++
		return h.storedCommands[h.index]
	}

	return ""
}

// AddCommand stores a previously run command in the shell model for future reuse. Only applied to commands which pass
// validation
func (h *History) AddCommand(command string) {
	if command == "" {
		return
	}

	binaryPathElements := strings.Split(command, "/")
	if len(binaryPathElements) == 0 {
		return
	}

	commandString := binaryPathElements[len(binaryPathElements)-1]
	h.storedCommands = append(h.storedCommands, commandString)
	h.index = len(h.storedCommands) - 1
}
