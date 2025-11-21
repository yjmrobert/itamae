package itamae

import (
	"github.com/charmbracelet/lipgloss"
)

// Tokyo Night color palette
// Reference: https://github.com/tokyo-night/tokyo-night-vscode-theme
var (
	// Background colors
	TokyoNightBg          = lipgloss.Color("#1a1b26") // Main background
	TokyoNightBgDark      = lipgloss.Color("#16161e") // Darker background
	TokyoNightBgFloat     = lipgloss.Color("#1f2335") // Floating elements
	TokyoNightBgHighlight = lipgloss.Color("#292e42") // Selection/highlight

	// Foreground colors
	TokyoNightFg      = lipgloss.Color("#c0caf5") // Main foreground
	TokyoNightFgDark  = lipgloss.Color("#a9b1d6") // Darker foreground
	TokyoNightComment = lipgloss.Color("#565f89") // Comments

	// Accent colors
	TokyoNightBlue    = lipgloss.Color("#7aa2f7") // Info/running
	TokyoNightCyan    = lipgloss.Color("#7dcfff") // Hints
	TokyoNightGreen   = lipgloss.Color("#9ece6a") // Success
	TokyoNightYellow  = lipgloss.Color("#e0af68") // Warning/skipped
	TokyoNightOrange  = lipgloss.Color("#ff9e64") // Warning alt
	TokyoNightRed     = lipgloss.Color("#f7768e") // Error
	TokyoNightMagenta = lipgloss.Color("#bb9af7") // Special
	TokyoNightPurple  = lipgloss.Color("#9d7cd8") // Special alt
)

// Status colors
var (
	StatusPending = TokyoNightComment
	StatusRunning = TokyoNightBlue
	StatusSuccess = TokyoNightGreen
	StatusError   = TokyoNightRed
	StatusSkipped = TokyoNightYellow
)

// Base styles
var (
	// Base style for all text
	BaseStyle = lipgloss.NewStyle().
			Foreground(TokyoNightFg).
			Background(TokyoNightBg)

	// Title style for pane headers
	TitleStyle = lipgloss.NewStyle().
			Foreground(TokyoNightMagenta).
			Background(TokyoNightBg).
			Bold(true).
			Padding(0, 1)

	// Border style
	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(TokyoNightBgHighlight).
			Background(TokyoNightBg)
)

// ChecklistStyle returns styled text for checklist items based on status
func ChecklistStyle(status string) lipgloss.Style {
	var color lipgloss.Color
	var icon string

	switch status {
	case "pending":
		color = StatusPending
		icon = "○"
	case "running":
		color = StatusRunning
		icon = "◉"
	case "success":
		color = StatusSuccess
		icon = "✓"
	case "error":
		color = StatusError
		icon = "✗"
	case "skipped":
		color = StatusSkipped
		icon = "⊘"
	default:
		color = TokyoNightComment
		icon = "?"
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Background(TokyoNightBg).
		SetString(icon)
}

// ChecklistPaneStyle returns the style for the checklist pane
func ChecklistPaneStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(TokyoNightBgDark).
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(TokyoNightBgHighlight)
}

// LogPaneStyle returns the style for the log pane
func LogPaneStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(TokyoNightBg).
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(TokyoNightBgHighlight)
}

// LogLineStyle returns styled text for log lines based on level
func LogLineStyle(level string) lipgloss.Style {
	var color lipgloss.Color

	switch level {
	case "info":
		color = TokyoNightFg
	case "success":
		color = StatusSuccess
	case "warning":
		color = StatusSkipped
	case "error":
		color = StatusError
	case "debug":
		color = TokyoNightComment
	default:
		color = TokyoNightFgDark
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Background(TokyoNightBg)
}

// ErrorPaneStyle returns the style for the error pane
func ErrorPaneStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(TokyoNightBgFloat).
		Padding(0, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(StatusError).
		BorderTop(true)
}

// ErrorTextStyle returns the style for error text
func ErrorTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(StatusError).
		Background(TokyoNightBgFloat).
		Bold(true)
}

// PhaseStyle returns styled text for installation phases
func PhaseStyle(phase string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(TokyoNightCyan).
		Background(TokyoNightBg).
		Bold(true)
}

// TimestampStyle returns style for log timestamps
func TimestampStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(TokyoNightComment).
		Background(TokyoNightBg)
}

// PackageNameStyle returns style for package names
func PackageNameStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(TokyoNightBlue).
		Background(TokyoNightBg)
}

// ProgressBarStyle returns style for progress indicators
func ProgressBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(TokyoNightMagenta).
		Background(TokyoNightBg)
}

// SummaryStyle returns style for the final summary
func SummaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(TokyoNightFg).
		Background(TokyoNightBgFloat).
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(TokyoNightMagenta).
		Bold(true)
}
