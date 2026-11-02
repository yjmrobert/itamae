# Adding Plugins

Learn how to add new tool plugins to Itamae.

## Plugin Structure

Each plugin is a shell script with:
1. Metadata block
2. `install()` function
3. `post_install()` function (optional)
4. `remove()` function
5. `check()` function
6. Router case statement

## Directory Selection

:::warning Choose Carefully
The directory determines how your plugin will be presented to users.
:::

Choose the appropriate directory:

- **`scripts/core/`**: Essential development tools (always installed together)
  - Examples: git, curl, nodejs, python3
  - Users cannot deselect individual packages
  
- **`scripts/essentials/`**: Common developer extras (installed as a group)
  - Examples: bat, ripgrep, zoxide, starship
  - All packages install together when category selected
  
- **`scripts/unverified/`**: Optional tools (user selectable)
  - Examples: vscode, zsh, ansible
  - Users can pick individual packages via multi-select

## Metadata Block

Start your script with metadata:

```bash
#!/bin/bash
#
# METADATA
# NAME: Tool Name
# DESCRIPTION: What the tool does
# INSTALL_METHOD: apt|binary|manual
# PACKAGE_NAME: actual-package-name  # For apt only
# REPO_SETUP: setup_repo              # Optional
# POST_INSTALL: post_install          # Optional
# REQUIRES: VAR_NAME|Prompt text|default command   # Optional, repeatable
#
```

### Metadata Fields

| Field | Required | Description |
|-------|----------|-------------|
| `NAME` | Yes | Display name of the tool |
| `DESCRIPTION` | Yes | Short description |
| `INSTALL_METHOD` | Yes | `apt`, `binary`, or `manual` |
| `PACKAGE_NAME` | For APT | Actual package name |
| `REPO_SETUP` | No | Function to add custom repository |
| `POST_INSTALL` | No | Function to run after installation |
| `REQUIRES` | No | User input required |

## Installation Methods

### APT Plugins

For tools available in APT repositories:

```bash
#!/bin/bash
#
# METADATA
# NAME: My Tool
# DESCRIPTION: A cool development tool
# INSTALL_METHOD: apt
# PACKAGE_NAME: my-tool
#

install() {
    echo "Installing My Tool..."
    if command -v nala &> /dev/null; then
        sudo nala install -y my-tool
    else
        sudo apt-get install -y my-tool
    fi
    echo "✅ My Tool installed."
}

remove() {
    echo "Removing My Tool..."
    sudo apt-get purge -y my-tool
    echo "✅ My Tool removed."
}

check() {
    command -v my-tool &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
```

### APT with Custom Repository

For tools requiring repository setup:

```bash
#!/bin/bash
#
# METADATA
# NAME: GitHub CLI
# DESCRIPTION: GitHub's official command line tool
# INSTALL_METHOD: apt
# PACKAGE_NAME: gh
# REPO_SETUP: setup_repo
#

setup_repo() {
    echo "Setting up GitHub CLI repository..."
    
    # Add GPG key
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | \
        sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
    
    # Add repository
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | \
        sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
    
    echo "✅ GitHub CLI repository added."
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
    sudo rm -f /usr/share/keyrings/githubcli-archive-keyring.gpg
    echo "✅ GitHub CLI removed."
}

check() {
    command -v gh &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    setup_repo) setup_repo ;;
    *) echo "Usage: $0 {install|remove|check|setup_repo}" && exit 1 ;;
esac
```

### APT with Symlink

For packages with different binary names:

```bash
#!/bin/bash
#
# METADATA
# NAME: bat (batcat)
# DESCRIPTION: A 'cat' clone with syntax highlighting
# INSTALL_METHOD: apt
# PACKAGE_NAME: bat
# POST_INSTALL: post_install
#

post_install() {
    # Create symlink: bat -> batcat
    mkdir -p "$HOME/.local/bin"
    ln -sf "$(command -v batcat)" "$HOME/.local/bin/bat"
    echo "✅ Created symlink: bat -> batcat"
}

install() {
    echo "Installing bat..."
    if command -v nala &> /dev/null; then
        sudo nala install -y bat
    else
        sudo apt-get install -y bat
    fi
    post_install
}

remove() {
    echo "Removing bat..."
    sudo apt-get purge -y bat
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

### Binary Plugins

For tools installed from binary releases:

```bash
#!/bin/bash
#
# METADATA
# NAME: Custom Binary
# DESCRIPTION: A tool distributed as binary
# INSTALL_METHOD: binary
#

install() {
    echo "Installing Custom Binary..."
    
    # Download binary
    wget https://example.com/custom-binary -O /tmp/custom-binary
    
    # Install
    chmod +x /tmp/custom-binary
    sudo mv /tmp/custom-binary /usr/local/bin/
    
    echo "✅ Custom Binary installed."
}

remove() {
    echo "Removing Custom Binary..."
    sudo rm -f /usr/local/bin/custom-binary
    echo "✅ Custom Binary removed."
}

check() {
    command -v custom-binary &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
```

## Batch Installation

Itamae optimizes APT installations through batching:

### Phase 0: Repository Setup
- All `REPO_SETUP` functions are called
- Single `apt-get update` runs after all repos added

### Phase 1: Batch APT Installation
- All APT packages installed in one command
- Post-install tasks run individually after batch completes

### Phase 2: Individual Installation
- Binary and manual plugins install one at a time

## Best Practices

- **Prefer `nala`, fall back to `apt-get`** in `install()` for APT tools; the batch phase does the same automatically.
- **Use `curl` in `setup_repo()`**. `curl`, `gnupg` and `ca-certificates` are installed before repository setup runs, `wget` is not.
- **Detect the architecture** with `uname -m` when downloading binaries. Releases ship for `amd64` and `arm64`; map to the upstream naming (`x86_64`/`aarch64`, `amd64`/`arm64`, `x64`/`arm64`) as needed.
- **Install user binaries into `~/.local/bin`**. Itamae puts `~/.local/bin`, `~/.cargo/bin` and `~/.atuin/bin` on `PATH` when running scripts and `check()`.
- **Install your own dependencies.** If the script needs `unzip`, `zip` or `fontconfig`, install them inside the script (see `cascadia-code.sh` or `sdkman.sh`).
- **Make `check()` honest.** It must exit 0 only when the tool is usable, and it must fail after `remove()`. Itamae runs it after every install; a tool whose check fails is reported as a failure, and the integration tests assert the removal side.
- **Keep scripts quiet about prompts.** Scripts run with `DEBIAN_FRONTEND=noninteractive` and no stdin; never ask questions. Declare inputs with `REQUIRES` instead.
- **Pass ShellCheck.** `shellcheck -S warning itamae/scripts/*/*.sh` runs in CI.

## Verification

After `install()` (or the batch install plus `post_install()`), Itamae runs `check()`. In the TUI the package moves through `Installing...`, `Running post-install...` and `Verifying...` before it is marked ✓. Manual plugins (`INSTALL_METHOD: manual`) are reported as "manual steps required" and are not verified.

## Testing Your Plugin

1. Register the plugin in `itamae/main_test.go`: add it to `expectedMetadata` (category, method, package name) and to `pluginAssertions` (a command fragment expected during `install` and `remove`). The suite fails for unregistered scripts, so the catalogue is always deliberate.
2. Run the mocked script tests:
   ```bash
   go test ./itamae -run 'TestInstall|TestUninstall|TestRepoSetup|TestMetadataParsing|TestRouterRejectsUnknownAction'
   ```
3. Try it for real, headlessly:
   ```bash
   go run . install --category unverified --plugins my-tool --yes --plain
   go run . check --plugins my-tool
   bash itamae/scripts/unverified/my-tool.sh remove
   ```
4. Push: the Integration workflow installs the whole category on Ubuntu 22.04/24.04 and Debian 12/13, with both `nala` and `apt-get`, and removes it again.

See [Testing](/developers/testing) for the full picture.

## Next Steps

- [Testing](/developers/testing) - Test your plugin
- [Contributing](/contributing) - Submit your plugin
