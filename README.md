# Itamae

Itamae is a command-line tool written in Go that sets up a developer's Linux workstation. It uses a plugin-based architecture to install and manage software.

## Installation

### Quick Install (Latest Version)

The easiest way to install Itamae is with the following command, which will download and run the install script:
```bash
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

### Install Specific Version

You can install a specific version by setting the `ITAMAE_VERSION` environment variable:
```bash
ITAMAE_VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

### Install with Go

If you have Go installed, you can install directly:
```bash
# Latest version
go install github.com/yjmrobert/itamae@latest

# Specific version
go install github.com/yjmrobert/itamae@v1.0.0
```
This will install the `itamae` binary in your `$HOME/go/bin` directory. Make sure that this directory is in your `PATH` to run the command from anywhere on your system.

### Build from Source

For development or customization:
```bash
git clone https://github.com/yjmrobert/itamae.git
cd itamae
./build.sh
sudo mv bin/itamae /usr/local/bin/
```

## Check Version

To check which version of Itamae you have installed:
```bash
itamae version
# or
itamae --version
```

## Commands

Itamae has the following commands:

*   **install:** Run `itamae install` to install a custom set of software.
*   **version:** Run `itamae version` or `itamae --version` to display version information.

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

### Commit Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated version management:

- `feat:` - New feature (minor version bump)
- `fix:` - Bug fix (patch version bump)  
- `feat!:` or `BREAKING CHANGE:` - Breaking change (major version bump)
- `docs:`, `chore:`, `style:`, `refactor:`, `test:` - Other changes (patch version bump)

Example:
```bash
git commit -m "feat: add support for Neovim plugin"
git commit -m "fix: resolve path issue in zsh script"
```

### Creating a Release

See [RELEASE.md](RELEASE.md) for detailed release instructions. Quick version:

```bash
./release.sh
```

This will automatically analyze commits, determine the version bump, run tests, and create a release.
