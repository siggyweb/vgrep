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

	logger, cleanUp := logging.ConfigureLogging()
	defer cleanUp()

	program := tea.NewProgram(vgrep.InitialModel(logger))
	if _, err := program.Run(); err != nil {
		fmt.Println(`an error occurred, exiting...`, err)
		os.Exit(1)
	}
}
