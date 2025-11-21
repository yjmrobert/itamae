package itamae

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderFullView renders the complete TUI with all panes
func renderFullView(m InstallModel) string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Render main panes (checklist + logs)
	checklistPane := renderChecklistPane(m)
	logPane := renderLogPane(m)

	// Join checklist and log panes horizontally
	mainView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		checklistPane,
		logPane,
	)

	// Add error pane at bottom if there are errors
	if len(m.errors) > 0 {
		errorPane := renderErrorPane(m)
		return lipgloss.JoinVertical(
			lipgloss.Left,
			mainView,
			errorPane,
		)
	}

	return mainView
}

// renderChecklistPane renders the left checklist pane
func renderChecklistPane(m InstallModel) string {
	if m.width == 0 {
		return ""
	}

	checklistWidth := int(float64(m.width) * 0.30)

	// Build header
	header := TitleStyle.Render("📦 Packages")
	phaseInfo := PhaseStyle(m.activePhase).Render(fmt.Sprintf("Phase: %s", m.activePhase))

	// Build package list
	var items []string
	items = append(items, header)
	items = append(items, phaseInfo)
	items = append(items, "") // Spacing

	for _, pkg := range m.packages {
		icon := ChecklistStyle(pkg.Status).Render("")

		// Format package line
		var line string
		if pkg.Status == "running" {
			// Show spinner for active package
			line = fmt.Sprintf("%s %s %s",
				m.spinner.View(),
				PackageNameStyle().Render(pkg.Name),
				lipgloss.NewStyle().Foreground(TokyoNightComment).Render(pkg.Progress),
			)
		} else {
			line = fmt.Sprintf("%s %s",
				icon,
				PackageNameStyle().Render(pkg.Name),
			)

			// Add progress/error on next line if present
			if pkg.Status == "error" && pkg.Error != "" {
				line += "\n  " + ErrorTextStyle().Width(checklistWidth-4).Render("↳ "+pkg.Error)
			} else if pkg.Progress != "" && pkg.Status != "pending" {
				line += "\n  " + lipgloss.NewStyle().
					Foreground(TokyoNightComment).
					Render("↳ "+pkg.Progress)
			}
		}

		items = append(items, line)
	}

	// Add summary if complete
	if m.complete {
		items = append(items, "")
		items = append(items, strings.Repeat("─", checklistWidth-4))
		items = append(items, fmt.Sprintf("✓ Success: %d", len(m.successful)))
		items = append(items, fmt.Sprintf("✗ Failed:  %d", len(m.failed)))
	}

	content := strings.Join(items, "\n")

	// Set viewport content and render
	m.checklistViewport.SetContent(content)

	viewportContent := m.checklistViewport.View()

	// Wrap in styled pane
	errorPaneHeight := 0
	if len(m.errors) > 0 {
		errorPaneHeight = 7
	}
	return ChecklistPaneStyle().
		Width(checklistWidth).
		Height(m.height - errorPaneHeight).
		Render(viewportContent)
}

// renderLogPane renders the right log pane
func renderLogPane(m InstallModel) string {
	if m.width == 0 {
		return ""
	}

	checklistWidth := int(float64(m.width) * 0.30)
	logWidth := m.width - checklistWidth - 6 // Account for borders and spacing

	// Build header
	header := TitleStyle.Render("📋 Installation Log")

	var items []string
	items = append(items, header)
	items = append(items, "") // Spacing

	// Add log lines
	for _, log := range m.logs {
		timestamp := TimestampStyle().Render(log.Timestamp.Format("15:04:05"))

		var prefix string
		switch log.Level {
		case "success":
			prefix = "✓"
		case "error":
			prefix = "✗"
		case "warning":
			prefix = "⚠"
		case "info":
			prefix = "ℹ"
		default:
			prefix = "·"
		}

		styledPrefix := LogLineStyle(log.Level).Render(prefix)

		message := log.Message
		if log.Package != "" {
			// Find package name
			pkgName := log.Package
			if idx, ok := m.packageIndex[log.Package]; ok {
				pkgName = m.packages[idx].Name
			}
			message = PackageNameStyle().Render(pkgName) + " → " + message
		}

		line := fmt.Sprintf("%s %s %s",
			timestamp,
			styledPrefix,
			LogLineStyle(log.Level).Render(message),
		)

		items = append(items, line)
	}

	// Show helpful hints if no logs yet
	if len(m.logs) == 0 {
		items = append(items, lipgloss.NewStyle().
			Foreground(TokyoNightComment).
			Render("Waiting for installation to begin..."))
	}

	// Add navigation hint at bottom if not complete
	if !m.complete {
		items = append(items, "")
		items = append(items, lipgloss.NewStyle().
			Foreground(TokyoNightComment).
			Italic(true).
			Render("Use ↑/↓ or j/k to scroll"))
	}

	content := strings.Join(items, "\n")

	// Set viewport content
	m.logViewport.SetContent(content)

	// Auto-scroll to bottom on new content
	if len(m.logs) > 0 {
		m.logViewport.GotoBottom()
	}

	viewportContent := m.logViewport.View()

	// Wrap in styled pane
	errorPaneHeight := 0
	if len(m.errors) > 0 {
		errorPaneHeight = 7
	}
	return LogPaneStyle().
		Width(logWidth).
		Height(m.height - errorPaneHeight).
		Render(viewportContent)
}

// renderErrorPane renders the bottom error pane
func renderErrorPane(m InstallModel) string {
	if len(m.errors) == 0 {
		return ""
	}

	// Build header
	header := ErrorTextStyle().Render(fmt.Sprintf("⚠ %d Error(s) Occurred", len(m.errors)))

	var items []string
	items = append(items, header)
	items = append(items, "") // Spacing

	// Show last 3 errors (most recent first)
	start := 0
	if len(m.errors) > 3 {
		start = len(m.errors) - 3
	}

	for i := len(m.errors) - 1; i >= start; i-- {
		err := m.errors[i]

		// Find package name
		pkgName := err.Package
		if idx, ok := m.packageIndex[err.Package]; ok {
			pkgName = m.packages[idx].Name
		}

		line := fmt.Sprintf("• %s [%s]: %s",
			PackageNameStyle().Render(pkgName),
			PhaseStyle(err.Phase).Render(err.Phase),
			err.Message,
		)

		items = append(items, ErrorTextStyle().Render(line))
	}

	content := strings.Join(items, "\n")

	// Wrap in styled pane
	return ErrorPaneStyle().
		Width(m.width - 4).
		Render(content)
}

// renderSummary renders the final installation summary
func renderSummary(m InstallModel) string {
	var items []string

	items = append(items, strings.Repeat("═", 60))
	items = append(items, "📊 INSTALLATION SUMMARY")
	items = append(items, strings.Repeat("═", 60))
	items = append(items, "")

	if len(m.successful) > 0 {
		items = append(items, lipgloss.NewStyle().
			Foreground(StatusSuccess).
			Bold(true).
			Render(fmt.Sprintf("✅ Successfully installed (%d):", len(m.successful))))

		for _, name := range m.successful {
			items = append(items, "   • "+name)
		}
		items = append(items, "")
	}

	if len(m.failed) > 0 {
		items = append(items, lipgloss.NewStyle().
			Foreground(StatusError).
			Bold(true).
			Render(fmt.Sprintf("❌ Failed to install (%d):", len(m.failed))))

		for _, name := range m.failed {
			items = append(items, "   • "+name)
		}
		items = append(items, "")
	}

	if len(m.failed) == 0 {
		items = append(items, lipgloss.NewStyle().
			Foreground(StatusSuccess).
			Bold(true).
			Render("🎉 All packages installed successfully!"))
	}

	items = append(items, "")
	items = append(items, strings.Repeat("═", 60))
	items = append(items, "Press 'q' or Ctrl+C to exit")

	return SummaryStyle().Render(strings.Join(items, "\n"))
}
