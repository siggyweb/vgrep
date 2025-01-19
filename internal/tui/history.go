package tui

type HistoryModel interface {
	GetPreviousCommand() string
	GetNextCommand() string
	AddCommand(string)
}

type History struct {
	storedCommands []string
	index          int
}

func (h *History) GetPreviousCommand() string {
	if len(h.storedCommands) == 0 {
		return ""
	} else if h.index > 0 {
		h.index--
		return h.storedCommands[h.index]
	}
	return ""
}

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

func (h *History) AddCommand(command string) {
	if command == "" {
		return
	}
	h.storedCommands = append(h.storedCommands, command)
	h.index = len(h.storedCommands) - 1
}
