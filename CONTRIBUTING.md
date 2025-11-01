# Contributing to Itamae

Contributions are welcome! Please follow these guidelines to contribute.

## Adding a New Plugin

To add a new plugin, create a new shell script in the `scripts/` directory. The script must have the following structure:

1.  **Metadata:** The script must start with a metadata block that provides information about the plugin.
    *   `# NAME:` The name of the software.
    *   `# OMAKASE:` Whether the plugin should be included in the default "omakase" installation (`true` or `false`).
    *   `# DESCRIPTION:` A short description of the software.

2.  **`install()` function:** This function should contain the commands to install the software.

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
  # Add installation commands here
}

remove() {
  echo "Removing My Awesome Tool..."
  # Add removal commands here
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

## Running the Linter and Tests

Before submitting a pull request, please run the linter and tests to ensure that your changes are correct and that you have not introduced any regressions.

```bash
go fmt ./...
go vet ./...
go test ./...
```
