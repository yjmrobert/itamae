# Developer Guide

This guide provides instructions for developers who want to add new tools to Itamae.

## How It Works

Itamae uses a Go binary to orchestrate shell scripts located in the `scripts/` directory. Each script is a plugin that knows how to install and remove a specific piece of software. The Go binary embeds these scripts and uses a terminal user interface (TUI) to let the user choose which plugins to run.

## Adding a New Plugin

To add a new plugin, create a new shell script in the `scripts/` directory. The script must have the following structure:

1.  **Metadata:** The script must start with a metadata block that provides information about the plugin.
    *   `# NAME:` The name of the software.
    *   `# OMAKASE:` Whether the plugin should be included in the default "omakase" installation (`true` or `false`).
    *   `# DESCRIPTION:` A short description of the software.

2.  **`install()` function:** This function should contain the commands to install the software.
    *   **Package Manager:** The `install()` function should prefer `nala` over `apt-get` if it is available.
    *   **Symlinks:** If a Debian/Ubuntu package uses a different binary name (e.g., `batcat`, `fd-find`), the script should create a symlink to the more common alias (e.g., `bat`, `fd`) in `$HOME/.local/bin`.

3.  **`remove()` function:** This function should contain the commands to remove the software.

4.  **Router:** The script must end with a `case` statement that calls the `install()` or `remove()` function based on the first argument (`$1`) passed to the script.

### Example

```bash
#!/bin/bash

# NAME: My Awesome Tool
# OMAKASE: true
# DESCRIPTION: A really cool tool that does awesome things.

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

case "$1" in
  install)
    install
    ;;
  remove)
    remove
    ;;
  *)
    echo "Usage: $0 {install|remove}"
    exit 1
    ;;
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
