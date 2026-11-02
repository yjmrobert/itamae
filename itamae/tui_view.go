package itamae

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Layout constants
const (
	checklistRatio = 0.30
	maxErrorsShown = 3
)

// paneWidths returns the outer widths of the checklist and log panes.
func (m InstallModel) paneWidths() (int, int) {
	checklistWidth := int(float64(m.width) * checklistRatio)
	if checklistWidth < 20 {
		checklistWidth = min(20, m.width/2)
	}
	logWidth := m.width - checklistWidth
	return checklistWidth, logWidth
}

// bottomPaneHeight returns the height of the error or summary pane shown below the main panes.
func (m InstallModel) bottomPaneHeight() int {
	if m.complete {
		return lipgloss.Height(renderSummary(m))
	}
	if len(m.errors) > 0 {
		return lipgloss.Height(renderErrorPane(m))
	}
	return 0
}

// layout resizes the viewports for the current window size.
func (m *InstallModel) layout() {
	if m.width == 0 || m.height == 0 {
		return
	}
	checklistWidth, logWidth := m.paneWidths()
	mainHeight := m.height - m.bottomPaneHeight()

	checklistStyle := ChecklistPaneStyle()
	logStyle := LogPaneStyle()

	m.checklistViewport.SetWidth(max(1, checklistWidth-checklistStyle.GetHorizontalFrameSize()))
	m.checklistViewport.SetHeight(max(1, mainHeight-checklistStyle.GetVerticalFrameSize()))
	m.logViewport.SetWidth(max(1, logWidth-logStyle.GetHorizontalFrameSize()))
	m.logViewport.SetHeight(max(1, mainHeight-logStyle.GetVerticalFrameSize()))
	m.syncViewports()
}

// syncViewports refreshes the viewport contents, keeping the log pinned to the
// bottom while the user has not scrolled away from it.
func (m *InstallModel) syncViewports() {
	m.checklistViewport.SetContent(renderChecklist(*m))

	atBottom := m.logViewport.AtBottom()
	m.logViewport.SetContent(renderLogLines(*m))
	if atBottom {
		m.logViewport.GotoBottom()
	}
}

// renderFullView renders the complete TUI with all panes
func renderFullView(m InstallModel) string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	checklistWidth, logWidth := m.paneWidths()
	mainHeight := m.height - m.bottomPaneHeight()

	checklistPane := ChecklistPaneStyle().
		Width(checklistWidth - ChecklistPaneStyle().GetHorizontalBorderSize()).
		Height(mainHeight - ChecklistPaneStyle().GetVerticalBorderSize()).
		Render(m.checklistViewport.View())
	logPane := LogPaneStyle().
		Width(logWidth - LogPaneStyle().GetHorizontalBorderSize()).
		Height(mainHeight - LogPaneStyle().GetVerticalBorderSize()).
		Render(m.logViewport.View())

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, checklistPane, logPane)

	switch {
	case m.complete:
		return lipgloss.JoinVertical(lipgloss.Left, mainView, renderSummary(m))
	case len(m.errors) > 0:
		return lipgloss.JoinVertical(lipgloss.Left, mainView, renderErrorPane(m))
	default:
		return mainView
	}
}

// renderChecklist renders the content of the left checklist pane
func renderChecklist(m InstallModel) string {
	width := m.checklistViewport.Width()

	items := []string{
		TitleStyle.Render("📦 Packages"),
		PhaseStyle().Render(fmt.Sprintf("Phase: %s", m.activePhase)),
		"",
	}

	for _, pkg := range m.packages {
		var line string
		if pkg.Status == "running" {
			line = fmt.Sprintf("%s %s", m.spinner.View(), PackageNameStyle().Render(pkg.Name))
		} else {
			line = fmt.Sprintf("%s %s", ChecklistStyle(pkg.Status).String(), PackageNameStyle().Render(pkg.Name))
		}

		detail := pkg.Progress
		detailStyle := lipgloss.NewStyle().Foreground(TokyoNightComment)
		if pkg.Status == "error" && pkg.Error != "" {
			detail = pkg.Error
			detailStyle = ErrorTextStyle().Background(TokyoNightBgDark)
		}
		if detail != "" && pkg.Status != "pending" {
			line += "\n  " + detailStyle.Width(max(1, width-2)).Render("↳ "+detail)
		}
		items = append(items, line)
	}

	return strings.Join(items, "\n")
}

// renderLogLines renders the content of the right log pane
func renderLogLines(m InstallModel) string {
	items := []string{TitleStyle.Render("📋 Installation Log"), ""}

	for _, log := range m.logs {
		timestamp := TimestampStyle().Render(log.Timestamp.Format("15:04:05"))
		prefix := LogLineStyle(log.Level).Render(LevelIcon(log.Level))

		message := log.Message
		if log.Package != "" {
			pkgName := log.Package
			if idx, ok := m.packageIndex[log.Package]; ok {
				pkgName = m.packages[idx].Name
			}
			message = PackageNameStyle().Render(pkgName) + " → " + message
		}

		items = append(items, fmt.Sprintf("%s %s %s", timestamp, prefix, LogLineStyle(log.Level).Render(message)))
	}

	if len(m.logs) == 0 {
		items = append(items, lipgloss.NewStyle().Foreground(TokyoNightComment).Render("Waiting for installation to begin..."))
	}

	if !m.complete {
		items = append(items, "", lipgloss.NewStyle().
			Foreground(TokyoNightComment).
			Italic(true).
			Render("Use ↑/↓ or j/k to scroll"))
	}

	return strings.Join(items, "\n")
}

// renderErrorPane renders the bottom error pane
func renderErrorPane(m InstallModel) string {
	if len(m.errors) == 0 {
		return ""
	}

	items := []string{ErrorTextStyle().Render(fmt.Sprintf("⚠ %d Error(s) Occurred", len(m.errors))), ""}

	start := 0
	if len(m.errors) > maxErrorsShown {
		start = len(m.errors) - maxErrorsShown
	}

	for i := len(m.errors) - 1; i >= start; i-- {
		err := m.errors[i]
		pkgName := err.Package
		if idx, ok := m.packageIndex[err.Package]; ok {
			pkgName = m.packages[idx].Name
		}
		if pkgName == "" {
			pkgName = "system"
		}
		items = append(items, ErrorTextStyle().Render(fmt.Sprintf("• %s [%s]: %s", pkgName, err.Phase, firstLine(err.Message))))
	}

	return ErrorPaneStyle().
		Width(max(1, m.width-ErrorPaneStyle().GetHorizontalBorderSize())).
		Render(strings.Join(items, "\n"))
}

// renderSummary renders the final installation summary shown once everything finished
func renderSummary(m InstallModel) string {
	var parts []string

	success := lipgloss.NewStyle().Foreground(StatusSuccess).Bold(true)
	failure := lipgloss.NewStyle().Foreground(StatusError).Bold(true)
	skipped := lipgloss.NewStyle().Foreground(StatusSkipped).Bold(true)

	parts = append(parts, success.Render(fmt.Sprintf("✓ %d installed", len(m.successful))))
	if len(m.failed) > 0 {
		parts = append(parts, failure.Render(fmt.Sprintf("✗ %d failed: %s", len(m.failed), strings.Join(m.failed, ", "))))
	}
	if len(m.skipped) > 0 {
		parts = append(parts, skipped.Render(fmt.Sprintf("⊘ %d need manual steps: %s", len(m.skipped), strings.Join(m.skipped, ", "))))
	}
	if len(m.failed) == 0 {
		parts = append(parts, success.Render("🎉 All packages installed successfully!"))
	}
	parts = append(parts, lipgloss.NewStyle().Foreground(TokyoNightComment).Render("Press q or Enter to exit"))

	return SummaryStyle().
		Width(max(1, m.width-SummaryStyle().GetHorizontalBorderSize())).
		Render(strings.Join(parts, "\n"))
}
