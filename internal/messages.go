package vgrep

import "time"

// CommandResponseMessage is the representation of a shell command that has been executed successfully or otherwise
type CommandResponseMessage struct {
	result string // could be []string depending on how result is composed? requires testing
	err    error
}

// TickMsg is used as a token to trigger running shell commands async which feed back into the view
type TickMsg time.Time
