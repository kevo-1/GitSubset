package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateModeSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "up" || msg.String() == "k":
			if m.modeCursor > 0 {
				m.modeCursor--
			}

		case msg.String() == "down" || msg.String() == "j":
			if m.modeCursor < len(m.modeChoices)-1 {
				m.modeCursor++
			}

		case msg.String() == "enter":
			switch FetchMode(m.modeCursor) {
			case ModeWholeRepo:
				m.selectedFiles = m.files
				m.screen = ScreenFetching
				return m, tea.Batch(m.spinner.Tick, fetchAllCmd(m.link.Path, m.files))

			case ModeFolders:
				m.picker = NewPickerModel(m.files, true)
				m.screen = ScreenPicker
				return m, nil

			case ModeFiles:
				m.picker = NewPickerModel(m.files, false)
				m.screen = ScreenPicker
				return m, nil
			}

		case msg.String() == "esc":
			m.screen = ScreenInput
			m.textInput.SetValue("")
			m.textInput.Focus()
			return m, nil

		case msg.String() == "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

// ---------- view ----------

func (m Model) viewModeSelect() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Select Fetch Mode"))
	b.WriteString("\n\n")

	repoInfo := fmt.Sprintf("Repository: %s/%s  •  %s",
		m.link.User, m.link.Repo, pluralize(len(m.files), "file", "files"))
	b.WriteString(infoBoxStyle.Render(repoInfo))
	b.WriteString("\n\n")

	b.WriteString(normalStyle.Render("What would you like to fetch?"))
	b.WriteString("\n\n")

	descriptions := []string{
		"Download every file in the repository",
		"Browse and select entire folders to download",
		"Browse and select individual files to download",
	}

	icons := []string{"[≡]", "[▸]", "[·]"}

	for i, choice := range m.modeChoices {
		cursor := "  "
		style := inactiveItemStyle
		if i == m.modeCursor {
			cursor = "▸ "
			style = activeItemStyle
		}

		line := fmt.Sprintf("%s%s %s", cursor, icons[i], style.Render(choice))
		b.WriteString(line)
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("    %s", dimStyle.Render(descriptions[i])))
		b.WriteString("\n\n")
	}

	return b.String()
}
