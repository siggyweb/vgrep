package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	vgrep "github.com/siggyweb/vgrep/internal"
	"os"
)

func main() {
	// fmt.Println("Welcome vgrep")
	program := tea.NewProgram(vgrep.InitialModel())
	if _, err := program.Run(); err != nil {
		fmt.Println(`an error occurred`, err)
		os.Exit(1)
	}
}
