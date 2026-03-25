package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#A78BFA")
	accentColor    = lipgloss.Color("#10B981")
	errorColor     = lipgloss.Color("#EF4444")
	warningColor   = lipgloss.Color("#F59E0B")
	dimColor       = lipgloss.Color("#6B7280")
	textColor      = lipgloss.Color("#F9FAFB")
	bgColor        = lipgloss.Color("#1F2937")

	// Title / Header
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 2).
			MarginBottom(1)

	// Subtitle
	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			MarginBottom(1)

	// Normal text
	normalStyle = lipgloss.NewStyle().
			Foreground(textColor)

	// Dim / help text
	dimStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	// Success
	successStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	// Error box
	errorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Foreground(errorColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Selected item
	selectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Active / highlighted item
	activeItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Bold(true).
			Padding(0, 1)

	// Inactive menu item
	inactiveItemStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Padding(0, 1)

	// Checked item
	checkedStyle = lipgloss.NewStyle().
			Foreground(accentColor)

	// Unchecked item
	uncheckedStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	// Footer / help bar
	helpBarStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			MarginTop(1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("#374151"))

	// Info box
	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2).
			MarginTop(1)

	// Spinner text
	spinnerTextStyle = lipgloss.NewStyle().
				Foreground(secondaryColor)
)
