---
title: Developer Guide
weight: 100
---

# Developer Guide

Welcome to the Itamae developer documentation. Here you'll learn how to contribute to Itamae and extend it with new plugins.

## Architecture

Itamae uses:
- **Go**: Main binary and orchestration
- **Shell Scripts**: Plugin-based tool installation
- **Bubbletea**: Terminal UI framework
- **Cobra**: CLI command structure
- **Hugo**: Documentation site

### Directory Structure

```
itamae/
├── cmd/              # CLI commands (install, logs, version)
├── itamae/           # Core logic and TUI
│   ├── scripts/      # Embedded shell scripts
│   │   ├── core/     # Essential tools
│   │   ├── essentials/ # Common extras
│   │   └── unverified/ # Optional tools
│   ├── tui_*.go      # Terminal UI components
│   └── debug_logger.go # Logging system
├── docs/             # Hugo documentation site
└── .github/workflows/ # CI/CD automation
```

## Quick Links

- [Adding Plugins]({{< relref "adding-plugins" >}}) - Create new tool plugins
- [Testing]({{< relref "testing" >}}) - Running and writing tests
- [TUI Development]({{< relref "tui" >}}) - Working with the terminal UI
- [Release Process]({{< relref "releases" >}}) - Creating new releases

## Contributing

See the [Contributing Guide]({{< relref "../contributing" >}}) for guidelines on:
- Code style
- Commit conventions
- Pull request process
