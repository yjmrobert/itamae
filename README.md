# Itamae

Itamae is a command-line tool written in Go that sets up a developer's Linux workstation. It uses a plugin-based architecture to install and manage software.

## Installation

The easiest way to install Itamae is with the following command, which will download and run the bootstrap script:
```bash
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/main/bootstrap.sh | bash
```

Alternatively, if you have Go installed, you can build and install it from source:
```bash
go install github.com/yjmrobert/itamae@latest
```
This will install the `itamae` binary in your `$HOME/go/bin` directory. Make sure that this directory is in your `PATH` to run the command from anywhere on your system.

## Commands

Itamae has two commands:

*   **install:** Run `itamae install` to install the core set of software and then choose from a list of additional plugins.
*   **uninstall:** Run `itamae uninstall` to remove all installed software.

## Available Plugins

The following plugins are available:

*   **Ripgrep:** A powerful command-line search tool.
*   **Visual Studio Code:** A popular code editor.

## How It Works

Itamae uses a Go binary to orchestrate shell scripts located in the `scripts/` directory. Each script is a plugin that knows how to install and remove a specific piece of software. The Go binary embeds these scripts and uses a terminal user interface (TUI) to let the user choose which plugins to run.

## Under the Hood

The TUI is built using the [Charm Bubble Tea](https://github.com/charmbracelet/bubbletea) library. The shell scripts are embedded into the Go binary using `go:embed`. The CLI is built using [Cobra](https://github.com/spf13/cobra).

## Contributing

Contributions are welcome! Please see the [Developer Guide](DEVELOPERS.md) for instructions on how to add a new tool.
