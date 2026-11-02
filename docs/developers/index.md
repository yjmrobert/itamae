# Developer Guide

Welcome to the Itamae developer documentation. Here you'll learn how to contribute to Itamae and extend it with new plugins.

## Architecture

Itamae uses:
- **Go 1.26**: Main binary and orchestration
- **Shell Scripts**: Plugin-based tool installation (embedded into the binary)
- **Charm** (`charm.land`): Bubble Tea v2 for the TUI, Bubbles v2 components, Lip Gloss v2 styling, Huh v2 forms, Log v2, and Fang for the CLI chrome (styled help, `--version`, man pages, completions)
- **Cobra**: CLI command structure
- **VitePress**: Documentation site

### Directory Structure

```
itamae/
├── cmd/                # CLI commands (install, check, list, logs, version)
├── itamae/             # Core logic and TUI
│   ├── plugin.go       # Metadata parsing and loading of embedded scripts
│   ├── engine.go       # Installer: phases, batching, fallback, verification
│   ├── exec.go         # Running scripts and package manager commands
│   ├── run.go          # TUI / plain front ends and summary output
│   ├── forms.go        # Huh forms
│   ├── tui_*.go        # Bubble Tea model and view
│   ├── theme.go        # Tokyo Night palette
│   ├── debug_logger.go # Logging system
│   └── scripts/        # Embedded shell scripts
│       ├── core/       # Essential tools
│       ├── essentials/ # Common extras
│       └── unverified/ # Optional tools
├── test/
│   ├── integration/    # Container harness used by the Integration workflow
│   └── install-sh/     # End-to-end tests for install.sh
├── docs/               # VitePress documentation site
└── .github/workflows/  # CI, Integration, Release, docs deployment
```

### Installation engine

`itamae/engine.go` contains the `Installer`, which is independent of any UI. It reports progress through the `Sender` interface:

- `*tea.Program` implements `Sender`, which drives the TUI
- `PlainSender` prints line-oriented output for `--plain`, CI and logs
- tests use a recording sender together with fake scripts and a mocked `sudo`

Every installed tool is verified with its `check()` function before it is reported as a success, and a failing batch APT install falls back to installing packages one by one.

## Quick Links

- [Adding Plugins](/developers/adding-plugins) - Create new tool plugins
- [Testing](/developers/testing) - Running and writing tests
- [TUI Development](/developers/tui) - Working with the terminal UI
- [Release Process](/developers/releases) - Creating new releases

## Contributing

See the [Contributing Guide](/contributing) for guidelines on:
- Code style
- Commit conventions
- Pull request process
