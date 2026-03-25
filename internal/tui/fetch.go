package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ---------- update ----------

func (m Model) updateDone(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			// Go back to mode select to pick more files
			m.screen = ScreenModeSelect
			m.modeCursor = 0
			return m, nil
		}
	}
	return m, nil
}

// ---------- view ----------

func (m Model) viewDone() string {
	var b strings.Builder

	b.WriteString(successStyle.Render("✓ Fetch Complete!"))
	b.WriteString("\n\n")

	details := fmt.Sprintf(
		"Repository:  %s/%s\nLocation:    %s\nFiles:       %s",
		m.link.User, m.link.Repo,
		m.link.Path,
		pluralize(m.fetchedCount, "file", "files"),
	)
	b.WriteString(infoBoxStyle.Render(details))
	b.WriteString("\n\n")

	b.WriteString(dimStyle.Render("Your files have been downloaded to the repository folder."))

	return b.String()
}
