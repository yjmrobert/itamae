package itamae

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// PackageStatus represents the current state of a package installation
type PackageStatus struct {
	ID        string    // Plugin ID
	Name      string    // Display name
	Status    string    // "pending", "running", "success", "error", "skipped"
	Progress  string    // Progress message (e.g., "Installing...", "Configuring...")
	Error     string    // Error message if status is "error"
	StartTime time.Time // When installation started
	EndTime   time.Time // When installation completed
}

// LogLine represents a single line in the log output
type LogLine struct {
	Timestamp time.Time
	Level     string // "info", "success", "warning", "error", "debug"
	Package   string // Which package this log is for (empty for system logs)
	Message   string
}

// ErrorInfo represents a structured error
type ErrorInfo struct {
	Package   string
	Phase     string // "repo_setup", "install", "post_install", "verify"
	Message   string
	Timestamp time.Time
}

// PhaseStartMsg indicates a new installation phase is starting
type PhaseStartMsg struct {
	Phase string // "repo_setup", "apt_batch", "individual"
	Count int    // Number of items in this phase
}

// PhaseCompleteMsg indicates a phase has completed
type PhaseCompleteMsg struct {
	Phase string
}

// PackageStartMsg indicates work on a package is starting
type PackageStartMsg struct {
	PackageID string
	Phase     string // "repo_setup", "install", "post_install", "verify"
}

// PackageCompleteMsg indicates a package installation completed
type PackageCompleteMsg struct {
	PackageID string
	Success   bool
	Skipped   bool   // The tool needs manual steps; neither success nor failure
	Error     string // Empty if Success is true (or a note when Skipped)
}

// LogMsg adds a new log line
type LogMsg struct {
	Level   string
	Package string // Empty for system logs
	Message string
}

// ErrorMsg adds a new error
type ErrorMsg struct {
	Package string
	Phase   string
	Message string
}

// ProgressMsg updates progress text for a package
type ProgressMsg struct {
	PackageID string
	Progress  string
}

// SummaryMsg signals installation is complete and shows summary
type SummaryMsg struct {
	Successful []string
	Failed     []string
	Skipped    []string
}

// InstallModel is the Bubble Tea model for the installation TUI
type InstallModel struct {
	packages     []PackageStatus
	packageIndex map[string]int // ID -> index mapping

	logs   []LogLine
	errors []ErrorInfo

	activePhase string // "init", "repo_setup", "apt_batch", "individual", "complete"

	checklistViewport viewport.Model
	logViewport       viewport.Model
	spinner           spinner.Model

	width  int
	height int

	successful []string
	failed     []string
	skipped    []string

	quitting bool
	complete bool
}

// NewInstallModel creates a new installation TUI model
func NewInstallModel(plugins []ToolPlugin) InstallModel {
	packages := make([]PackageStatus, len(plugins))
	packageIndex := make(map[string]int, len(plugins))
	for i, p := range plugins {
		packages[i] = PackageStatus{ID: p.ID, Name: p.Name, Status: "pending"}
		packageIndex[p.ID] = i
	}

	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(TokyoNightMagenta)),
	)

	m := InstallModel{
		packages:          packages,
		packageIndex:      packageIndex,
		activePhase:       "init",
		checklistViewport: viewport.New(),
		logViewport:       viewport.New(),
		spinner:           s,
	}
	m.syncViewports()
	return m
}

// Init initializes the model
func (m InstallModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles messages and updates the model
func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "q", "enter", "esc":
			if m.complete {
				m.quitting = true
				return m, tea.Quit
			}
		case "up", "k":
			m.logViewport.ScrollUp(1)
		case "down", "j":
			m.logViewport.ScrollDown(1)
		case "pgup":
			m.logViewport.PageUp()
		case "pgdown":
			m.logViewport.PageDown()
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout()
		return m, nil

	case PhaseStartMsg:
		m.activePhase = msg.Phase
		m.addLog("info", "", fmt.Sprintf("Starting phase: %s (%d items)", msg.Phase, msg.Count))

	case PhaseCompleteMsg:
		m.addLog("success", "", fmt.Sprintf("Completed phase: %s", msg.Phase))

	case PackageStartMsg:
		if idx, ok := m.packageIndex[msg.PackageID]; ok {
			m.packages[idx].Status = "running"
			m.packages[idx].Progress = phaseLabel(msg.Phase)
			if m.packages[idx].StartTime.IsZero() {
				m.packages[idx].StartTime = time.Now()
			}
		}

	case ProgressMsg:
		if idx, ok := m.packageIndex[msg.PackageID]; ok {
			m.packages[idx].Progress = msg.Progress
		}

	case PackageCompleteMsg:
		if idx, ok := m.packageIndex[msg.PackageID]; ok {
			m.packages[idx].EndTime = time.Now()
			switch {
			case msg.Skipped:
				m.packages[idx].Status = "skipped"
				m.packages[idx].Progress = msg.Error
			case msg.Success:
				m.packages[idx].Status = "success"
				m.packages[idx].Progress = "Complete"
			default:
				m.packages[idx].Status = "error"
				m.packages[idx].Error = msg.Error
			}
		}

	case LogMsg:
		m.addLog(msg.Level, msg.Package, msg.Message)

	case ErrorMsg:
		m.errors = append(m.errors, ErrorInfo{
			Package:   msg.Package,
			Phase:     msg.Phase,
			Message:   msg.Message,
			Timestamp: time.Now(),
		})
		m.layout()

	case SummaryMsg:
		m.activePhase = "complete"
		m.complete = true
		m.successful = msg.Successful
		m.failed = msg.Failed
		m.skipped = msg.Skipped
		m.layout()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		m.syncViewports()
		return m, cmd

	default:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		m.checklistViewport, cmd = m.checklistViewport.Update(msg)
		cmds = append(cmds, cmd)
		m.logViewport, cmd = m.logViewport.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	m.syncViewports()
	return m, nil
}

func phaseLabel(phase string) string {
	switch phase {
	case "repo_setup":
		return "Configuring repository..."
	case "install":
		return "Installing..."
	case "post_install":
		return "Running post-install..."
	case "verify":
		return "Verifying..."
	default:
		return phase
	}
}

// addLog is a helper to add a log line
func (m *InstallModel) addLog(level, pkg, message string) {
	m.logs = append(m.logs, LogLine{
		Timestamp: time.Now(),
		Level:     level,
		Package:   pkg,
		Message:   message,
	})
}

// View renders the TUI
func (m InstallModel) View() tea.View {
	var content string
	if m.quitting && !m.complete {
		content = "Installation cancelled.\n"
	} else {
		content = renderFullView(m)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
