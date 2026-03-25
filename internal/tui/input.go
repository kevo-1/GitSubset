package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ---------- update ----------

func (m Model) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			url := strings.TrimSpace(m.textInput.Value())
			if url == "" {
				m.inputErr = "Please enter a repository URL"
				return m, nil
			}
			if !strings.Contains(url, "github.com") {
				m.inputErr = "Only GitHub URLs are supported (must contain github.com)"
				return m, nil
			}
			m.inputErr = ""
			m.screen = ScreenCloning
			return m, tea.Batch(m.spinner.Tick, cloneCmd(url))

		case "q":
			if m.textInput.Value() == "" {
				return m, tea.Quit
			}
		case "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// ---------- view ----------

func (m Model) viewInput() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Enter a GitHub repository URL"))
	b.WriteString("\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n")

	if m.inputErr != "" {
		b.WriteString("\n")
		b.WriteString(errorBoxStyle.Render("⚠  " + m.inputErr))
	}

	return b.String()
}
