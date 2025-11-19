---
title: Developer Guide
weight: 4
---

# Developer Guide

This guide provides instructions for developers who want to add new tools to Itamae.

## How It Works

Itamae uses a Go binary to orchestrate shell scripts located in the `scripts/` directory. Each script is a plugin that knows how to install and remove a specific piece of software. The Go binary embeds these scripts and uses interactive forms (built with [Charm Huh](https://github.com/charmbracelet/huh)) to let the user choose which optional plugins to install.

## Adding a New Plugin

To add a new plugin, create a new shell script in the `scripts/` directory. The script must have the following structure:

1.  **Metadata:** The script must start with a metadata block that provides information about the plugin.
    *   `# NAME:` The name of the software.
    *   `# OMAKASE:` Whether the plugin should be included in the default "omakase" installation (`true` or `false`).
    *   `# DESCRIPTION:` A short description of the software.
    *   `# INSTALL_METHOD:` The installation method: `apt`, `binary`, or `manual`.
    *   `# PACKAGE_NAME:` (For `apt` plugins only) The actual package name in the APT repository.
    *   `# POST_INSTALL:` (Optional) The name of a function to run after batch APT installation (e.g., `post_install`).

2.  **`install()` function:** This function should contain the commands to install the software.
    *   **Package Manager:** For APT-based tools, the `install()` function should prefer `nala` over `apt-get` if it is available. This function serves as a fallback for individual installations.
    *   **Symlinks:** If a Debian/Ubuntu package uses a different binary name (e.g., `batcat`, `fd-find`), create a `post_install()` function that creates symlinks to the more common alias (e.g., `bat`, `fd`) in `$HOME/.local/bin`.

3.  **`post_install()` function:** (Optional, for APT plugins) This function runs after batch installation to perform post-installation tasks like creating symlinks or configuration.

4.  **`remove()` function:** This function should contain the commands to remove the software.

5.  **`check()` function:** This function should return 0 if the software is installed, non-zero otherwise.

6.  **Router:** The script must end with a `case` statement that calls the appropriate function based on the first argument (`$1`) passed to the script.

### Batch Installation Optimization

Itamae optimizes APT package installation by batching all APT-based tools into a single `nala install` or `apt-get install` command. This significantly improves performance by:
- Resolving all dependencies in a single pass
- Reducing package manager overhead
- Providing cleaner output with a single progress indicator

Tools are automatically categorized by their `INSTALL_METHOD`:
- **`apt`**: Batch installed with nala/apt-get, then post-install tasks are run individually
- **`binary`**: Installed individually using their custom installation scripts
- **`manual`**: Require manual installation by the user

### Example: APT Plugin with Symlink

```bash
#!/bin/bash
#
# METADATA
# NAME: bat (batcat)
# OMAKASE: true
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
# OMAKASE: true
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
# OMAKASE: false
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

## Testing

Unit tests for shell script plugins are written in Go (`main_test.go`). This approach was chosen as a compromise to full integration tests due to environment limitations. A `TestMain` function sets up a mock environment that replaces system commands (e.g., `sudo`, `apt-get`) with logging scripts. Tests then assert that the specific, correct install/remove commands are called by verifying the log file's contents.

Before submitting a pull request, please run the linter and tests to ensure that your changes are correct and that you have not introduced any regressions.

```bash
go fmt ./...
go vet ./...
go test ./...
```
