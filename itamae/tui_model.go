package itamae

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	Phase     string // "repo_setup", "install", "post_install"
	Message   string
	Timestamp time.Time
}

// InstallModel is the Bubbletea model for the installation TUI
type InstallModel struct {
	// Package tracking
	packages     []PackageStatus
	packageIndex map[string]int // ID -> index mapping

	// Logs and errors
	logs   []LogLine
	errors []ErrorInfo

	// Current state
	activePhase string // "init", "repo_setup", "apt_batch", "individual", "summary", "complete"
	currentPkg  string // ID of currently installing package

	// UI components
	checklistViewport viewport.Model
	logViewport       viewport.Model
	spinner           spinner.Model

	// Dimensions
	width  int
	height int

	// Installation result
	successful []string // Package names
	failed     []string // Package names

	// Control
	quitting bool
	complete bool
}

// Message types for Bubbletea updates

// WindowSizeMsg is sent when terminal is resized (handled by tea.WindowSizeMsg)

// PhaseStartMsg indicates a new installation phase is starting
type PhaseStartMsg struct {
	Phase string // "repo_setup", "apt_batch", "individual", "summary"
	Count int    // Number of items in this phase
}

// PhaseCompleteMsg indicates a phase has completed
type PhaseCompleteMsg struct {
	Phase string
}

// PackageStartMsg indicates a package installation is starting
type PackageStartMsg struct {
	PackageID string
	Phase     string // "repo_setup", "install", "post_install"
}

// PackageCompleteMsg indicates a package installation completed
type PackageCompleteMsg struct {
	PackageID string
	Success   bool
	Error     string // Empty if Success is true
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

// ProgressMsg updates progress text for current package
type ProgressMsg struct {
	PackageID string
	Progress  string
}

// SummaryMsg signals installation is complete and shows summary
type SummaryMsg struct {
	Successful []string
	Failed     []string
}

// SpinnerTickMsg is sent by the spinner
type SpinnerTickMsg = spinner.TickMsg

// NewInstallModel creates a new installation TUI model
func NewInstallModel(plugins []ToolPlugin) InstallModel {
	// Initialize packages list
	packages := make([]PackageStatus, len(plugins))
	packageIndex := make(map[string]int)

	for i, p := range plugins {
		packages[i] = PackageStatus{
			ID:     p.ID,
			Name:   p.Name,
			Status: "pending",
		}
		packageIndex[p.ID] = i
	}

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(TokyoNightMagenta)

	// Initialize viewports (will be sized properly in Init)
	checklistVP := viewport.New(0, 0)
	logVP := viewport.New(0, 0)

	return InstallModel{
		packages:          packages,
		packageIndex:      packageIndex,
		logs:              []LogLine{},
		errors:            []ErrorInfo{},
		activePhase:       "init",
		checklistViewport: checklistVP,
		logViewport:       logVP,
		spinner:           s,
		successful:        []string{},
		failed:            []string{},
	}
}

// Init initializes the model
func (m InstallModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.EnterAltScreen,
	)
}

// Update handles messages and updates the model
func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.complete {
				m.quitting = true
				return m, tea.Quit
			}
		case "up", "k":
			m.logViewport.LineUp(1)
		case "down", "j":
			m.logViewport.LineDown(1)
		case "pgup":
			m.logViewport.ViewUp()
		case "pgdown":
			m.logViewport.ViewDown()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate pane dimensions
		// Checklist: 30% width, full height minus error pane
		// Log: 70% width, full height minus error pane
		// Error: full width, 5 lines (if errors exist)

		errorHeight := 0
		if len(m.errors) > 0 {
			errorHeight = 5
		}

		mainHeight := m.height - errorHeight - 4 // -4 for borders and padding

		checklistWidth := int(float64(m.width) * 0.30)
		logWidth := m.width - checklistWidth - 4 // -4 for borders

		m.checklistViewport.Width = checklistWidth
		m.checklistViewport.Height = mainHeight

		m.logViewport.Width = logWidth
		m.logViewport.Height = mainHeight

	case PhaseStartMsg:
		m.activePhase = msg.Phase
		m.addLog("info", "", fmt.Sprintf("Starting phase: %s (%d items)", msg.Phase, msg.Count))

	case PhaseCompleteMsg:
		m.addLog("success", "", fmt.Sprintf("Completed phase: %s", msg.Phase))

	case PackageStartMsg:
		m.currentPkg = msg.PackageID
		if idx, ok := m.packageIndex[msg.PackageID]; ok {
			m.packages[idx].Status = "running"
			m.packages[idx].Progress = fmt.Sprintf("Phase: %s", msg.Phase)
			m.packages[idx].StartTime = time.Now()
		}
		m.addLog("info", msg.PackageID, fmt.Sprintf("Starting %s", msg.Phase))

	case PackageCompleteMsg:
		if idx, ok := m.packageIndex[msg.PackageID]; ok {
			m.packages[idx].EndTime = time.Now()
			if msg.Success {
				m.packages[idx].Status = "success"
				m.packages[idx].Progress = "Complete"
				m.successful = append(m.successful, m.packages[idx].Name)
				m.addLog("success", msg.PackageID, "Installation successful")
			} else {
				m.packages[idx].Status = "error"
				m.packages[idx].Error = msg.Error
				m.failed = append(m.failed, m.packages[idx].Name)
				m.addLog("error", msg.PackageID, fmt.Sprintf("Installation failed: %s", msg.Error))
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

	case ProgressMsg:
		if idx, ok := m.packageIndex[msg.PackageID]; ok {
			m.packages[idx].Progress = msg.Progress
		}

	case SummaryMsg:
		m.activePhase = "complete"
		m.complete = true
		m.successful = msg.Successful
		m.failed = msg.Failed
		return m, nil

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Update viewports
	m.checklistViewport, cmd = m.checklistViewport.Update(msg)
	cmds = append(cmds, cmd)

	m.logViewport, cmd = m.logViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// addLog is a helper to add a log line and update the viewport
func (m *InstallModel) addLog(level, pkg, message string) {
	m.logs = append(m.logs, LogLine{
		Timestamp: time.Now(),
		Level:     level,
		Package:   pkg,
		Message:   message,
	})

	// Auto-scroll log viewport to bottom
	m.logViewport.GotoBottom()
}

// View renders the TUI
func (m InstallModel) View() string {
	if m.quitting {
		return "Installation cancelled.\n"
	}

	// Delegate to tui_view.go for rendering
	return renderFullView(m)
}
