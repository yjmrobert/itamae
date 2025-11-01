# Itamae

Itamae is a command-line tool written in Go that sets up a developer's Linux workstation. It uses a plugin-based architecture to install and manage software.

## Modes

Itamae has three modes:

*   **Omakase (Default):** Run `go run .` to install a pre-selected set of essential software.
*   **Customize:** Run `go run . customize` to choose which software to install from a list of available plugins.
*   **Remove:** Run `go run . remove` to choose which software to remove.

## Available Plugins

The following plugins are available:

*   **Ripgrep:** A powerful command-line search tool.
*   **Visual Studio Code:** A popular code editor.

## How It Works

Itamae uses a Go binary to orchestrate shell scripts located in the `scripts/` directory. Each script is a plugin that knows how to install and remove a specific piece of software. The Go binary embeds these scripts and uses a terminal user interface (TUI) to let the user choose which plugins to run.

## Under the Hood

The TUI is built using the [Charm Bubble Tea](https://github.com/charmbracelet/bubbletea) library. The shell scripts are embedded into the Go binary using `go:embed`.
