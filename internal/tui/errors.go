package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ---------- update ----------

func (m Model) updateError(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			// Retry the failed operation
			switch m.errPrevScreen {
			case ScreenInput:
				url := strings.TrimSpace(m.textInput.Value())
				m.screen = ScreenCloning
				return m, tea.Batch(m.spinner.Tick, cloneCmd(url))

			case ScreenCloning:
				m.screen = ScreenListing
				return m, tea.Batch(m.spinner.Tick, listCmd(m.link.Path))

			case ScreenPicker:
				m.screen = ScreenFetching
				return m, tea.Batch(m.spinner.Tick, fetchCmd(m.link.Path, m.selectedFiles))
			}

		case "esc":
			// Go back to a safe screen
			switch m.errPrevScreen {
			case ScreenInput:
				m.screen = ScreenInput
				m.textInput.Focus()
			case ScreenCloning:
				m.screen = ScreenInput
				m.textInput.Focus()
			case ScreenPicker:
				m.screen = ScreenModeSelect
				m.modeCursor = 0
			default:
				m.screen = ScreenInput
				m.textInput.Focus()
			}
			return m, nil
		}
	}
	return m, nil
}

// ---------- view ----------

func (m Model) viewError() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("An Error Occurred"))
	b.WriteString("\n")
	b.WriteString(errorBoxStyle.Render("✕  " + m.errMsg))
	b.WriteString("\n\n")

	b.WriteString(normalStyle.Render("What would you like to do?"))
	b.WriteString("\n\n")
	b.WriteString(normalStyle.Render("  r  →  Retry the operation"))
	b.WriteString("\n")
	b.WriteString(normalStyle.Render("  esc →  Go back"))
	b.WriteString("\n")
	b.WriteString(normalStyle.Render("  q  →  Quit"))

	return b.String()
}
