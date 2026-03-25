package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ---------- update (cloning) ----------

func (m Model) updateCloning(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case cloneDoneMsg:
		m.link = msg.link
		m.clonedByUs = msg.clonedByUs
		m.screen = ScreenListing
		return m, tea.Batch(m.spinner.Tick, listCmd(m.link.Path))

	case cloneErrMsg:
		m.errMsg = fmt.Sprintf("Failed to clone repository:\n%s", msg.err.Error())
		m.errPrevScreen = ScreenInput
		m.screen = ScreenError
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// ---------- update (listing) ----------

func (m Model) updateListing(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case listDoneMsg:
		m.files = msg.files
		m.screen = ScreenModeSelect
		m.modeCursor = 0
		return m, nil

	case listErrMsg:
		m.errMsg = fmt.Sprintf("Failed to list repository contents:\n%s", msg.err.Error())
		m.errPrevScreen = ScreenCloning
		m.screen = ScreenError
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// ---------- update (fetching) ----------

func (m Model) updateFetching(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fetchDoneMsg:
		m.fetchedCount = msg.count
		m.fetched = true
		m.screen = ScreenDone
		return m, nil

	case fetchErrMsg:
		m.errMsg = fmt.Sprintf("Failed to fetch files:\n%s", msg.err.Error())
		m.errPrevScreen = ScreenPicker
		m.screen = ScreenError
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// ---------- views ----------

func (m Model) viewCloning() string {
	var b strings.Builder
	b.WriteString(subtitleStyle.Render("Cloning Repository"))
	b.WriteString("\n\n")
	b.WriteString(spinnerTextStyle.Render(m.spinner.View() + " Cloning repository (sparse, no-checkout)…"))
	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render("This may take a moment for large repositories."))
	return b.String()
}

func (m Model) viewListing() string {
	var b strings.Builder
	b.WriteString(subtitleStyle.Render("Loading Repository Structure"))
	b.WriteString("\n\n")
	b.WriteString(spinnerTextStyle.Render(m.spinner.View() + " Listing files in the repository…"))
	return b.String()
}

func (m Model) viewFetching() string {
	var b strings.Builder
	b.WriteString(subtitleStyle.Render("Fetching Files"))
	b.WriteString("\n\n")
	b.WriteString(spinnerTextStyle.Render(
		fmt.Sprintf("%s Fetching %s…", m.spinner.View(), pluralize(len(m.selectedFiles), "file", "files")),
	))
	return b.String()
}
