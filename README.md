# Itamae

Itamae is a command-line tool written in Go that sets up a developer's Linux workstation. It uses a plugin-based architecture to install and manage software.

## Installation

The easiest way to install Itamae is with the following command, which will download and run the install script:
```bash
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
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

Itamae uses a Go binary to orchestrate shell scripts located in the `scripts/` directory. Each script is a plugin that knows how to install and remove a specific piece of software. The Go binary embeds these scripts and provides an interactive form interface for selecting optional plugins.

### Installation Process

1. **Interactive Selection:** Choose which additional tools to install using an interactive multi-select form
2. **Installation Plan:** Review the complete installation plan organized by method (APT, binary, manual)
3. **Batch Installation:** All APT packages are installed in a single optimized command using `nala` (or `apt-get`)
4. **Individual Installation:** Binary and manual installations run individually with live progress

### Performance Optimization

Itamae optimizes installation by batching all APT packages into a single command, significantly improving speed and dependency resolution compared to installing packages one-by-one.

## Under the Hood

The interactive forms are built using the [Charm Huh](https://github.com/charmbracelet/huh) library. The shell scripts are embedded into the Go binary using `go:embed`. The CLI is built using [Cobra](https://github.com/spf13/cobra).

## Contributing

Contributions are welcome! Please see the [Developer Guide](DEVELOPERS.md) for instructions on how to add a new tool.
