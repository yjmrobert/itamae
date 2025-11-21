---
title: TUI Development
weight: 30
---

# TUI Development

Learn about Itamae's Terminal User Interface architecture and how to extend it.

## Overview

Itamae's TUI is built with the [Charm](https://charm.sh/) ecosystem, providing a modern, interactive terminal experience using:

- **[Bubbletea](https://github.com/charmbracelet/bubbletea)**: The Elm-inspired framework for building terminal apps
- **[Bubbles](https://github.com/charmbracelet/bubbles)**: Pre-built TUI components (viewport, spinner, etc.)
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)**: Style definitions and layouts
- **[Huh](https://github.com/charmbracelet/huh)**: Interactive forms for user input

## Architecture

### Component Structure

```
itamae/
├── theme.go           # Tokyo Night color scheme and styles
├── tui_model.go       # Bubbletea model and state management
├── tui_view.go        # View rendering (layout and display)
├── tui_orchestrator.go # Installation orchestration
└── executor.go        # Async command execution
```

### Data Flow

```
User Input → Huh Forms → Bubbletea Model → Orchestrator → Executor
                              ↓                              ↓
                         View Rendering ← Message Updates ← Script Output
```

## Theme System

### Tokyo Night Color Palette

Defined in `itamae/theme.go`:

```go
var (
    ColorBackground   = lipgloss.Color("#1a1b26")  // Dark background
    ColorForeground   = lipgloss.Color("#c0caf5")  // Light text
    ColorComment      = lipgloss.Color("#565f89")  // Dim gray
    ColorCyan         = lipgloss.Color("#7dcfff")  // Bright cyan
    ColorGreen        = lipgloss.Color("#9ece6a")  // Success green
    ColorYellow       = lipgloss.Color("#e0af68")  // Warning yellow
    ColorRed          = lipgloss.Color("#f7768e")  // Error red
    ColorPurple       = lipgloss.Color("#bb9af7")  // Accent purple
    ColorBlue         = lipgloss.Color("#7aa2f7")  // Info blue
)
```

### Style Definitions

Common styles:

```go
var (
    StyleBase        = lipgloss.NewStyle().Foreground(ColorForeground)
    StyleTitle       = StyleBase.Bold(true).Foreground(ColorCyan)
    StyleSuccess     = StyleBase.Foreground(ColorGreen)
    StyleError       = StyleBase.Foreground(ColorRed)
    StyleInfo        = StyleBase.Foreground(ColorBlue)
    StyleDim         = StyleBase.Foreground(ColorComment)
    StyleBorder      = lipgloss.NewStyle().BorderForeground(ColorComment)
)
```

### Using Styles

```go
// Apply style to text
styledText := StyleSuccess.Render("✅ Installation complete")

// Combine styles
titleStyle := StyleTitle.Bold(true).Underline(true)
```

## Bubbletea Model

### Model Structure

The main model in `itamae/tui_model.go`:

```go
type InstallModel struct {
    packages       []Package      // List of packages to install
    logs           []string       // Installation logs
    errors         []string       // Error messages
    viewport       viewport.Model // Scrollable log pane
    errorViewport  viewport.Model // Scrollable error pane
    spinner        spinner.Model  // Loading spinner
    width          int           // Terminal width
    height         int           // Terminal height
    installDone    bool          // Installation complete flag
    currentPackage string        // Currently installing package
}
```

### Message Types

Messages drive state updates:

```go
// Phase transition (repository setup, batch install, individual install)
type PhaseStartMsg struct {
    Phase int
    Name  string
}

// Package installation started
type PackageStartMsg struct {
    Name string
}

// Package installation completed
type PackageCompleteMsg struct {
    Name    string
    Success bool
    Error   string
}

// Log line received
type LogMsg struct {
    Line string
}

// Error occurred
type ErrorMsg struct {
    Error string
}

// Installation summary
type SummaryMsg struct {
    Successful int
    Failed     int
}
```

### Update Function

Handle messages and update state:

```go
func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case PhaseStartMsg:
        // Update current phase
        m.logs = append(m.logs, formatPhaseStart(msg))
        
    case PackageStartMsg:
        // Update current package
        m.currentPackage = msg.Name
        
    case PackageCompleteMsg:
        // Mark package as complete
        updatePackageStatus(m.packages, msg.Name, msg.Success)
        
    case LogMsg:
        // Append to logs
        m.logs = append(m.logs, msg.Line)
        m.viewport.SetContent(strings.Join(m.logs, "\n"))
        m.viewport.GotoBottom()
        
    case ErrorMsg:
        // Add error
        m.errors = append(m.errors, msg.Error)
        m.errorViewport.SetContent(strings.Join(m.errors, "\n"))
        
    case tea.KeyMsg:
        // Handle keyboard input
        switch msg.String() {
        case "q", "ctrl+c":
            if m.installDone {
                return m, tea.Quit
            }
        case "j", "down":
            m.viewport.LineDown(1)
        case "k", "up":
            m.viewport.LineUp(1)
        }
    }
    
    return m, nil
}
```

## View Rendering

### Layout System

The view in `itamae/tui_view.go` uses a three-pane layout:

```
┌────────────────────────────────────────────────────┐
│                    Header                          │
├──────────────────┬─────────────────────────────────┤
│   Checklist      │      Logs                       │
│   (30% width)    │      (70% width)                │
│                  │                                 │
│ ⏳ git           │ [12:34:56] Installing git...    │
│ ⏳ nodejs        │ [12:34:57] apt-get install git  │
│ ⏳ python        │ [12:34:58] ✅ git installed     │
│                  │                                 │
├──────────────────┴─────────────────────────────────┤
│                  Errors (if any)                   │
└────────────────────────────────────────────────────┘
```

### Checklist Rendering

```go
func (m InstallModel) renderChecklist() string {
    var items []string
    
    for _, pkg := range m.packages {
        icon := getStatusIcon(pkg.Status)
        style := getStatusStyle(pkg.Status)
        
        if pkg.Name == m.currentPackage {
            icon = m.spinner.View()
        }
        
        item := fmt.Sprintf("%s %s", icon, pkg.Name)
        items = append(items, style.Render(item))
    }
    
    return lipgloss.JoinVertical(lipgloss.Left, items...)
}
```

### Log Rendering

```go
func (m InstallModel) renderLogs() string {
    return StyleBorder.
        Border(lipgloss.RoundedBorder()).
        Width(int(float64(m.width) * 0.7)).
        Height(m.height - 10).
        Render(m.viewport.View())
}
```

### Error Rendering

```go
func (m InstallModel) renderErrors() string {
    if len(m.errors) == 0 {
        return ""
    }
    
    return StyleError.
        Border(lipgloss.RoundedBorder()).
        BorderForeground(ColorRed).
        Width(m.width - 4).
        Render(m.errorViewport.View())
}
```

## Async Execution

### Running Scripts

Execute shell scripts asynchronously in `itamae/executor.go`:

```go
func ExecuteScriptAsyncCmd(scriptPath string, args []string) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command(scriptPath, args...)
        
        // Capture output
        stdout, _ := cmd.StdoutPipe()
        stderr, _ := cmd.StderrPipe()
        
        // Start command
        if err := cmd.Start(); err != nil {
            return ErrorMsg{Error: err.Error()}
        }
        
        // Stream output
        go func() {
            scanner := bufio.NewScanner(stdout)
            for scanner.Scan() {
                // Send log message
                program.Send(LogMsg{Line: scanner.Text()})
            }
        }()
        
        // Wait for completion
        err := cmd.Wait()
        
        return PackageCompleteMsg{
            Name:    scriptName,
            Success: err == nil,
            Error:   errorString(err),
        }
    }
}
```

## Interactive Forms

### Package Selection Form

Using Huh for multi-select:

```go
func showPackageSelectionForm(packages []string) ([]string, error) {
    var selected []string
    
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewMultiSelect[string]().
                Title("Select packages to install").
                Options(huh.NewOptions(packages...)...).
                Value(&selected).
                Height(15),
        ),
    ).WithTheme(huh.ThemeCharm())
    
    if err := form.Run(); err != nil {
        return nil, err
    }
    
    return selected, nil
}
```

### Confirmation Dialog

```go
func showConfirmation(message string) (bool, error) {
    var confirm bool
    
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewConfirm().
                Title(message).
                Value(&confirm),
        ),
    ).WithTheme(huh.ThemeCharm())
    
    if err := form.Run(); err != nil {
        return false, err
    }
    
    return confirm, nil
}
```

## Orchestration

### Installation Flow

Coordinated in `itamae/tui_orchestrator.go`:

```go
func processInstallTUI(packages []Package, program *tea.Program) {
    // Phase 0: Repository Setup
    program.Send(PhaseStartMsg{Phase: 0, Name: "Repository Setup"})
    for _, pkg := range packages {
        if pkg.RepoSetup != "" {
            runRepoSetup(pkg, program)
        }
    }
    
    // Phase 1: Batch APT Installation
    program.Send(PhaseStartMsg{Phase: 1, Name: "Batch Installation"})
    aptPackages := filterAPTPackages(packages)
    runBatchInstall(aptPackages, program)
    
    // Phase 2: Individual Installation
    program.Send(PhaseStartMsg{Phase: 2, Name: "Individual Installation"})
    for _, pkg := range filterNonAPTPackages(packages) {
        runIndividualInstall(pkg, program)
    }
    
    // Summary
    program.Send(SummaryMsg{
        Successful: countSuccessful(packages),
        Failed:     countFailed(packages),
    })
}
```

## Best Practices

### Style Consistency

- Always use theme colors
- Apply consistent spacing and borders
- Use status icons consistently

### Performance

- Use viewports for large content
- Batch updates when possible
- Debounce rapid updates

### Error Handling

- Always display errors to users
- Log errors for debugging
- Provide recovery options

### Accessibility

- Use clear status indicators
- Provide keyboard navigation
- Display help text when needed

## Testing TUI

### Manual Testing

```bash
# Build and run
./build.sh
./bin/itamae install
```

### Debugging

Add debug output:

```go
// In Update function
fmt.Fprintf(os.Stderr, "DEBUG: Received message: %T\n", msg)
```

View debug output:

```bash
./bin/itamae install 2> debug.log
```

## Extending the TUI

### Adding New Panes

1. Update model with new viewport
2. Add rendering function in `tui_view.go`
3. Update layout to include new pane

### Adding New Messages

1. Define message type in `tui_model.go`
2. Handle in `Update()` function
3. Send from orchestrator or executor

### Customizing Theme

1. Modify colors in `theme.go`
2. Update style definitions
3. Rebuild and test

## Resources

- [Bubbletea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lipgloss Docs](https://github.com/charmbracelet/lipgloss)
- [Huh Examples](https://github.com/charmbracelet/huh/tree/main/examples)

## Next Steps

- [Testing]({{< relref "testing" >}}) - Test your TUI changes
- [Release Process]({{< relref "releases" >}}) - Prepare for release
