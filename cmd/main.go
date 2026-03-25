package main

import (
	"fmt"
	"os"

	"GitSubset/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running GitSubset: %v\n", err)
		os.Exit(1)
	}

	if model, ok := m.(tui.Model); ok {
		model.Cleanup()
	}
}
