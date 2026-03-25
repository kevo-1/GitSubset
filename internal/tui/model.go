package tui

import (
	"fmt"
	"os"
	"strings"

	"GitSubset/internal"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents the current TUI screen.
type Screen int

const (
	ScreenInput      Screen = iota
	ScreenCloning
	ScreenListing
	ScreenModeSelect
	ScreenPicker
	ScreenFetching
	ScreenDone
	ScreenError
)

// FetchMode represents what the user wants to fetch.
type FetchMode int

const (
	ModeWholeRepo FetchMode = iota
	ModeFolders
	ModeFiles
)

// ---------- messages ----------

type cloneDoneMsg struct {
	link       internal.GithubLink
	clonedByUs bool
}
type cloneErrMsg struct{ err error }
type listDoneMsg struct{ files []string }
type listErrMsg struct{ err error }
type fetchDoneMsg struct{ count int }
type fetchErrMsg struct{ err error }

// ---------- Model ----------

type Model struct {
	screen Screen

	// URL input
	textInput textinput.Model
	inputErr  string

	// Clone / fetch spinner
	spinner spinner.Model

	// Repo info
	link  internal.GithubLink
	files []string

	// Mode selection
	modeChoices []string
	modeCursor  int

	// Picker
	picker PickerModel

	// Fetch
	selectedFiles []string
	fetchedCount  int

	// Cleanup tracking
	clonedByUs bool
	fetched    bool

	// Error
	errMsg       string
	errPrevScreen Screen

	// Window size
	width  int
	height int
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "https://github.com/user/repo"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return Model{
		screen:      ScreenInput,
		textInput:   ti,
		spinner:     sp,
		modeChoices: []string{"Whole Repository", "Select Folders", "Select Files"},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global quit on ctrl+c
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.screen {
	case ScreenInput:
		return m.updateInput(msg)
	case ScreenCloning:
		return m.updateCloning(msg)
	case ScreenListing:
		return m.updateListing(msg)
	case ScreenModeSelect:
		return m.updateModeSelect(msg)
	case ScreenPicker:
		return m.updatePicker(msg)
	case ScreenFetching:
		return m.updateFetching(msg)
	case ScreenDone:
		return m.updateDone(msg)
	case ScreenError:
		return m.updateError(msg)
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render("  GitSubset  "))
	b.WriteString("\n\n")

	switch m.screen {
	case ScreenInput:
		b.WriteString(m.viewInput())
	case ScreenCloning:
		b.WriteString(m.viewCloning())
	case ScreenListing:
		b.WriteString(m.viewListing())
	case ScreenModeSelect:
		b.WriteString(m.viewModeSelect())
	case ScreenPicker:
		b.WriteString(m.viewPicker())
	case ScreenFetching:
		b.WriteString(m.viewFetching())
	case ScreenDone:
		b.WriteString(m.viewDone())
	case ScreenError:
		b.WriteString(m.viewError())
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(m.viewHelp())

	return b.String()
}

func (m Model) viewHelp() string {
	var help string
	switch m.screen {
	case ScreenInput:
		help = "enter: confirm • q/ctrl+c: quit"
	case ScreenCloning, ScreenListing, ScreenFetching:
		help = "ctrl+c: quit"
	case ScreenModeSelect:
		help = "↑/↓: navigate • enter: select • esc: back • q: quit"
	case ScreenPicker:
		help = "↑/↓: navigate • space: toggle • tab: expand/collapse • a: all • enter: confirm • esc: back"
	case ScreenDone:
		help = "r: select more • q: quit"
	case ScreenError:
		help = "r: retry • esc: back • q: quit"
	}
	return helpBarStyle.Render(dimStyle.Render(help))
}

// ---------- clone command ----------

func cloneCmd(url string) tea.Cmd {
	return func() tea.Msg {
		// Check if the repo directory already exists before cloning
		parsedLink, parseErr := internal.ParseURL(url)
		dirExisted := false
		if parseErr == nil {
			if _, statErr := os.Stat(parsedLink.Repo); statErr == nil {
				dirExisted = true
			}
		}

		link, err := internal.Clone(url)
		if err != nil {
			return cloneErrMsg{err}
		}
		return cloneDoneMsg{link: link, clonedByUs: !dirExisted}
	}
}

// ---------- list command ----------

func listCmd(repoPath string) tea.Cmd {
	return func() tea.Msg {
		files, err := internal.ListContent(repoPath)
		if err != nil {
			return listErrMsg{err}
		}
		return listDoneMsg{files}
	}
}

// ---------- fetch command ----------

func fetchCmd(repoPath string, files []string) tea.Cmd {
	return func() tea.Msg {
		if err := internal.FetchContent(repoPath, files); err != nil {
			return fetchErrMsg{err}
		}
		return fetchDoneMsg{len(files)}
	}
}

// ---------- fetch all (whole repo) ----------

func fetchAllCmd(repoPath string, files []string) tea.Cmd {
	return func() tea.Msg {
		if err := internal.FetchContent(repoPath, files); err != nil {
			return fetchErrMsg{err}
		}
		return fetchDoneMsg{len(files)}
	}
}

// ---------- helpers ----------

// Cleanup removes the cloned repo if we created it and no files were fetched.
func (m Model) Cleanup() {
	if m.clonedByUs && !m.fetched && m.link.Path != "" {
		os.RemoveAll(m.link.Path)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clamp(v, lo, hi int) int {
	return max(lo, min(v, hi))
}

func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %s", n, plural)
}
