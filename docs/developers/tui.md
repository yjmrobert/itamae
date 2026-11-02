# TUI Development

Learn about Itamae's Terminal User Interface architecture and how to extend it.

## Overview

Itamae's TUI is built with the [Charm](https://charm.sh/) ecosystem, using the v2 generation of the libraries published under the `charm.land` module path:

- **[Bubble Tea v2](https://github.com/charmbracelet/bubbletea)** (`charm.land/bubbletea/v2`): the Elm-inspired framework for terminal apps
- **[Bubbles v2](https://github.com/charmbracelet/bubbles)** (`charm.land/bubbles/v2`): viewport and spinner components
- **[Lip Gloss v2](https://github.com/charmbracelet/lipgloss)** (`charm.land/lipgloss/v2`): styles and layout
- **[Huh v2](https://github.com/charmbracelet/huh)** (`charm.land/huh/v2`): interactive forms
- **[Log v2](https://github.com/charmbracelet/log)** (`charm.land/log/v2`): structured logging
- **[Fang](https://github.com/charmbracelet/fang)**: styled help, errors, `--version`, man pages and completions for the Cobra CLI

## Architecture

### Component Structure

```
itamae/
├── theme.go       # Tokyo Night color scheme and styles
├── tui_model.go   # Bubble Tea model, messages and state management
├── tui_view.go    # View rendering (layout and display)
├── run.go         # RunInstallTUI / RunInstallPlain, PlainSender, PrintSummary
├── forms.go       # Huh forms (category, multi-select, confirm, text input)
└── engine.go      # Installer: produces the messages the TUI consumes
```

### Data Flow

```
Huh forms ──► cmd/install.go ──► Installer.Run(ctx)
                                     │  Sender.Send(msg)
                                     ▼
                     *tea.Program (TUI)  or  PlainSender (--plain / CI)
                                     │
                                     ▼
                           InstallModel.Update ──► View
```

The `Installer` never touches the UI directly. It emits messages through the `Sender` interface, so the same engine powers the TUI, the plain-text mode and the unit tests.

## Theme System

Colours are plain `color.Color` values created with `lipgloss.Color`, defined in `itamae/theme.go`:

```go
var (
    TokyoNightBg      = lipgloss.Color("#1a1b26")
    TokyoNightFg      = lipgloss.Color("#c0caf5")
    TokyoNightComment = lipgloss.Color("#565f89")
    TokyoNightBlue    = lipgloss.Color("#7aa2f7") // info / running
    TokyoNightGreen   = lipgloss.Color("#9ece6a") // success
    TokyoNightYellow  = lipgloss.Color("#e0af68") // warning / skipped
    TokyoNightRed     = lipgloss.Color("#f7768e") // error
    TokyoNightMagenta = lipgloss.Color("#bb9af7") // titles
)
```

Helpers such as `StatusColor(status)`, `StatusIcon(status)`, `LevelColor(level)` and `LevelIcon(level)` keep icons and colours consistent between the TUI, the plain output and `itamae logs`.

Lip Gloss v2 does not detect the terminal's colour support by itself. Bubble Tea handles that for the TUI; everything printed outside the TUI goes through a `colorprofile.Writer`, which downgrades colours to what the terminal supports and strips them when the output is not a terminal.

## Bubble Tea Model

### Messages

The installer sends these messages (defined in `tui_model.go`):

| Message | Meaning |
|---------|---------|
| `PhaseStartMsg{Phase, Count}` | A phase (`repo_setup`, `apt_batch`, `individual`) begins |
| `PhaseCompleteMsg{Phase}` | A phase finished |
| `PackageStartMsg{PackageID, Phase}` | Work on a package starts (`repo_setup`, `install`, `post_install`, `verify`) |
| `ProgressMsg{PackageID, Progress}` | Progress text for a package |
| `PackageCompleteMsg{PackageID, Success, Skipped, Error}` | Final state of a package |
| `LogMsg{Level, Package, Message}` | A log line |
| `ErrorMsg{Package, Phase, Message}` | A structured error for the error pane |
| `SummaryMsg{Successful, Failed, Skipped}` | The run is complete |

### Model and Update

```go
type InstallModel struct {
    packages          []PackageStatus
    logs              []LogLine
    errors            []ErrorInfo
    checklistViewport viewport.Model
    logViewport       viewport.Model
    spinner           spinner.Model
    width, height     int
    complete          bool
}

func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:          // v2: KeyPressMsg instead of KeyMsg
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "q", "enter", "esc":
            if m.complete {
                return m, tea.Quit
            }
        case "up", "k":
            m.logViewport.ScrollUp(1)
        case "down", "j":
            m.logViewport.ScrollDown(1)
        }
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        m.layout()
    case PackageCompleteMsg:
        // update the checklist entry
    }
    m.syncViewports()
    return m, nil
}
```

### View

In Bubble Tea v2 `View()` returns a `tea.View` rather than a string. Alternate screen and mouse support are requested on the view:

```go
func (m InstallModel) View() tea.View {
    v := tea.NewView(renderFullView(m))
    v.AltScreen = true
    v.MouseMode = tea.MouseModeCellMotion
    return v
}
```

## View Rendering

The layout has two main panes and an optional bottom pane:

```
┌──────────────────┬─────────────────────────────────┐
│ 📦 Packages      │ 📋 Installation Log             │
│ Phase: apt_batch │ 12:34:56 ℹ Starting phase ...   │
│ ✓ git            │ 12:34:57 ✓ git → Installed      │
│ ◉ nodejs         │ 12:34:58 ℹ nodejs → Verifying   │
│   ↳ Verifying... │                                 │
│ ○ helm           │                                 │
├──────────────────┴─────────────────────────────────┤
│ ⚠ 1 Error(s) Occurred   /   📊 summary when done   │
└────────────────────────────────────────────────────┘
```

- `renderChecklist` renders the left pane (30% of the width) with status icons and a spinner for running packages
- `renderLogLines` renders timestamped, colour-coded log lines; the viewport stays pinned to the bottom until the user scrolls up
- `renderErrorPane` shows the three most recent errors
- `renderSummary` replaces the error pane once the run is complete

Viewport sizes are derived from the pane styles (`GetHorizontalFrameSize`, `GetVerticalFrameSize`) so borders and padding never push content off screen.

## Plain Mode

`RunInstallPlain` uses `PlainSender`, which renders the same messages as indented, colour-aware lines, and streams the raw script output underneath. It is used for `--plain`, whenever stdout is not a terminal, and in CI.

## Interactive Forms

Forms live in `itamae/forms.go` and use Huh v2:

```go
func SelectPlugins(plugins []ToolPlugin) ([]ToolPlugin, error) {
    var selectedIDs []string
    form := huh.NewForm(huh.NewGroup(
        huh.NewMultiSelect[string]().
            Title("Available Tools").
            Options(options...).
            Value(&selectedIDs).
            Height(30),
    ))
    if err := form.Run(); err != nil {
        return nil, err
    }
    return FilterPlugins(plugins, selectedIDs)
}
```

`cmd/install.go` only shows a form when the corresponding flag was not given and both stdin and stdout are terminals.

## Testing the TUI

The engine is covered by unit tests with a recording `Sender`; the TUI itself is best checked manually:

```bash
./build.sh
./bin/itamae install                       # TUI
./bin/itamae install -c core --yes --plain # same engine, plain output
```

## Extending the TUI

### Adding New Panes

1. Add a viewport to `InstallModel` and size it in `layout()`
2. Add a rendering function in `tui_view.go`
3. Include it in `renderFullView`

### Adding New Messages

1. Define the message type in `tui_model.go`
2. Handle it in `Update()` and, if useful, in `PlainSender.Send`
3. Send it from `engine.go`

### Customizing the Theme

Modify the colours in `theme.go`; every style derives from them.

## Resources

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [Huh](https://github.com/charmbracelet/huh)
- [Fang](https://github.com/charmbracelet/fang)

## Next Steps

- [Testing](/developers/testing) - Test your changes
- [Release Process](/developers/releases) - Prepare for release
