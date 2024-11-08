package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/siggyweb/vgrep/internal/logging"
	vgrep "github.com/siggyweb/vgrep/internal/tui"
	"os"
)

func main() {
	//goland:noinspection ALL
	fmt.Println("Welcome to vgrep! The dynamic terminal wrapper")

	file := logging.ConfigureLogging()
	defer file()

	program := tea.NewProgram(vgrep.InitialModel())
	if _, err := program.Run(); err != nil {
		fmt.Println(`an error occurred, exiting...`, err)
		os.Exit(1)
	}
}
