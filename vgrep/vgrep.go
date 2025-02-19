package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/siggyweb/vgrep/internal/logging"
	"github.com/siggyweb/vgrep/internal/stats"
	vgrep "github.com/siggyweb/vgrep/internal/tui"
	"os"
)

func main() {
	fmt.Println("launching dynamic terminal...")

	logger, cleanUpLogger := logging.ConfigureLogging()
	statsModel := &stats.SessionStatsModel{}

	defer func() {
		logger.Infof(statsModel.GetSummary())
		_ = cleanUpLogger()
	}()

	program := tea.NewProgram(vgrep.InitialModel(logger, statsModel))
	if _, err := program.Run(); err != nil {
		fmt.Println(`an error occurred, exiting...`, err)
		os.Exit(1)
	}
}
