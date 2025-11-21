# Developer Guide

This guide provides instructions for developers who want to add new tools to Itamae.

## How It Works

Itamae uses a Go binary to orchestrate shell scripts organized into three directories:
- **`scripts/core/`** - Essential development tools that are always installed together
- **`scripts/essentials/`** - Common developer extras that are installed as a group
- **`scripts/unverified/`** - Additional tools that users can selectively install

Each script is a plugin that knows how to install and remove a specific piece of software. The Go binary embeds these scripts and provides a full-featured TUI (Terminal User Interface) built with the [Charm ecosystem](https://charm.sh/) (Bubbletea, Bubbles, Lipgloss, and Huh) that displays:
- **Left pane**: Checklist of packages with real-time status updates (pending, running, success, error)
- **Right pane**: Scrollable installation logs with color-coded output
- **Bottom pane**: Error messages when failures occur

The TUI uses the **Tokyo Night** color scheme for a modern, readable appearance.

## Adding a New Plugin

To add a new plugin, create a new shell script in either `scripts/core/` (for essential tools) or `scripts/unverified/` (for optional tools). The script must have the following structure:

1.  **Metadata:** The script must start with a metadata block that provides information about the plugin.
    *   `# NAME:` The name of the software.
    *   `# DESCRIPTION:` A short description of the software.
    *   `# INSTALL_METHOD:` The installation method: `apt`, `binary`, or `manual`.
    *   `# PACKAGE_NAME:` (For `apt` plugins only) The actual package name in the APT repository.
    *   `# REPO_SETUP:` (Optional, for APT plugins) The name of a function that adds custom repositories (e.g., `setup_repo`). This is called before the batch `apt-get update`.
    *   `# POST_INSTALL:` (Optional) The name of a function to run after batch APT installation (e.g., `post_install`).
    *   `# REQUIRES:` (Optional) Required user inputs in format `VAR_NAME|Prompt text` (can have multiple).

2.  **`install()` function:** This function should contain the commands to install the software.
    *   **Package Manager:** For APT-based tools, the `install()` function should prefer `nala` over `apt-get` if it is available. This function serves as a fallback for individual installations.
    *   **Symlinks:** If a Debian/Ubuntu package uses a different binary name (e.g., `batcat`, `fd-find`), create a `post_install()` function that creates symlinks to the more common alias (e.g., `bat`, `fd`) in `$HOME/.local/bin`.

3.  **`post_install()` function:** (Optional, for APT plugins) This function runs after batch installation to perform post-installation tasks like creating symlinks or configuration.

4.  **`remove()` function:** This function should contain the commands to remove the software.

5.  **`check()` function:** This function should return 0 if the software is installed, non-zero otherwise.

6.  **Router:** The script must end with a `case` statement that calls the appropriate function based on the first argument (`$1`) passed to the script.

### Batch Installation Optimization

Itamae optimizes APT package installation by batching all APT-based tools into a single `nala install` or `apt-get install` command. The installation process follows three phases:

**Phase 0: Repository Setup**
- All plugins with `REPO_SETUP` metadata have their `setup_repo()` functions called sequentially
- A single `nala update` or `apt-get update` is run after all repositories are added
- This enables packages from custom repositories (GitHub CLI, Node.js, .NET, Java) to be installed in the batch phase

**Phase 1: Batch APT Installation**
- All APT packages (both standard and from custom repos) are installed in a single parallel command
- Post-install tasks are run individually after batch installation completes
- This significantly improves performance by resolving all dependencies in one pass

**Phase 2: Individual Installation**
- Binary and manual plugins are installed one at a time using their custom scripts

Tools are automatically categorized by their `INSTALL_METHOD`:
- **`apt`**: Batch installed with nala/apt-get, leveraging parallel capabilities
- **`binary`**: Installed individually using their custom installation scripts
- **`manual`**: Require manual installation by the user

### Example: APT Plugin with Symlink

```bash
#!/bin/bash
#
# METADATA
# NAME: bat (batcat)
# DESCRIPTION: A 'cat' clone with syntax highlighting and Git integration.
# INSTALL_METHOD: apt
# PACKAGE_NAME: batcat
# POST_INSTALL: post_install
#

post_install() {
    # Create the 'bat' symlink that all tools expect
    mkdir -p "$HOME/.local/bin"
    ln -sf "$(command -v batcat)" "$HOME/.local/bin/bat"
    echo "✅ Created symlink: bat -> batcat"
}

install() {
    echo "Installing bat..."
    if command -v nala &> /dev/null; then
        sudo nala install -y batcat
    else
        sudo apt-get install -y batcat
    fi
    post_install
}

remove() {
    echo "Removing bat..."
    sudo apt-get purge -y batcat
    rm -f "$HOME/.local/bin/bat"
    echo "✅ bat removed."
}

check() {
    command -v bat &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    post_install) post_install ;;
    *) echo "Usage: $0 {install|remove|check|post_install}" && exit 1 ;;
esac
```

### Example: Simple APT Plugin

```bash
#!/bin/bash
#
# METADATA
# NAME: My Awesome Tool
# DESCRIPTION: A really cool tool that does awesome things.
# INSTALL_METHOD: apt
# PACKAGE_NAME: my-awesome-tool
#

install() {
    echo "Installing My Awesome Tool..."
    if command -v nala &> /dev/null; then
        sudo nala install -y my-awesome-tool
    else
        sudo apt-get install -y my-awesome-tool
    fi
    echo "✅ My Awesome Tool installed."
}

remove() {
    echo "Removing My Awesome Tool..."
    sudo apt-get purge -y my-awesome-tool
    echo "✅ My Awesome Tool removed."
}

check() {
    command -v my-awesome-tool &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
```

### Example: Binary Plugin

```bash
#!/bin/bash
#
# METADATA
# NAME: My Binary Tool
# DESCRIPTION: A tool installed from a binary release.
# INSTALL_METHOD: binary
#

install() {
    echo "Installing My Binary Tool..."
    curl -L "https://example.com/tool.tar.gz" | sudo tar -xz -C /usr/local/bin
    echo "✅ My Binary Tool installed."
}

remove() {
    echo "Removing My Binary Tool..."
    sudo rm /usr/local/bin/tool
    echo "✅ My Binary Tool removed."
}

check() {
    command -v tool &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
```

### Example: APT Plugin with Custom Repository

For packages that need custom repositories (like GitHub CLI, Node.js, .NET, Java), add a `setup_repo()` function:

```bash
#!/bin/bash
#
# METADATA
# NAME: GitHub CLI
# DESCRIPTION: Official GitHub command-line tool.
# INSTALL_METHOD: apt
# PACKAGE_NAME: gh
# REPO_SETUP: setup_repo
#

setup_repo() {
    echo "Setting up GitHub CLI repository..."
    sudo mkdir -p /etc/apt/keyrings
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
    echo "✅ GitHub CLI repository configured."
}

install() {
    echo "Installing GitHub CLI..."
    if command -v nala &> /dev/null; then
        sudo nala install -y gh
    else
        sudo apt-get install -y gh
    fi
    echo "✅ GitHub CLI installed."
}

remove() {
    echo "Removing GitHub CLI..."
    sudo apt-get purge -y gh
    sudo rm -f /etc/apt/sources.list.d/github-cli.list
    sudo rm -f /etc/apt/keyrings/githubcli-archive-keyring.gpg
    echo "✅ GitHub CLI removed."
}

check() {
    command -v gh &> /dev/null
}

# --- ROUTER ---
case "$1" in
    setup_repo) setup_repo ;;
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {setup_repo|install|remove|check}" && exit 1 ;;
esac
```

## TUI Architecture

The installation interface is built using the [Charm ecosystem](https://charm.sh/) with a multi-pane layout:

### Components

**Theme System** (`itamae/theme.go`)
- Tokyo Night color palette with Lipgloss styles
- Color-coded status indicators (pending: gray, running: blue, success: green, error: red)
- Consistent styling across all panes

**Bubbletea Model** (`itamae/tui_model.go`)
- `InstallModel`: Main TUI state containing packages, logs, errors, and UI components
- Message types: `PhaseStartMsg`, `PackageStartMsg`, `PackageCompleteMsg`, `LogMsg`, `ErrorMsg`, `SummaryMsg`
- Implements `tea.Model` interface with `Init()`, `Update()`, and `View()` methods

**View Rendering** (`itamae/tui_view.go`)
- `renderChecklist()`: Left pane (30% width) showing package status with spinners
- `renderLogs()`: Right pane (70% width) with scrollable, timestamped logs
- `renderErrors()`: Bottom pane (5 lines) displaying recent errors when present
- Bubbles viewports for smooth scrolling

**Async Execution** (`itamae/executor.go`)
- `ExecuteScriptAsyncCmd()`: Runs shell scripts and captures output
- `ExecuteBatchAPTCmd()`: Executes parallel APT installations
- Returns `tea.Cmd` to send messages back to TUI

**Orchestration** (`itamae/tui_orchestrator.go`)
- `RunInstallTUI()`: Entry point that initializes TUI and starts installation
- `processInstallTUI()`: Coordinates three-phase installation and sends progress updates

### User Interactions

- **↑/↓ or j/k**: Scroll log pane
- **PgUp/PgDown**: Page through logs
- **q or Ctrl+C**: Exit (only after installation completes)

The TUI automatically scrolls logs to show the latest output and highlights the currently installing package with a spinner.

## Testing

Unit tests for shell script plugins are written in Go (`main_test.go`). This approach was chosen as a compromise to full integration tests due to environment limitations. A `TestMain` function sets up a mock environment that replaces system commands (e.g., `sudo`, `apt-get`) with logging scripts. Tests then assert that the specific, correct install/remove commands are called by verifying the log file's contents.

Before submitting a pull request, please run the linter and tests to ensure that your changes are correct and that you have not introduced any regressions.

```bash
go fmt ./...
go vet ./...
go test ./...
```
