# Usage Guide

Learn how to use Itamae to set up your development environment.

## Commands

| Command | Description |
|---------|-------------|
| `itamae install` | Install a category of tools, interactively or headlessly |
| `itamae check` | Report which tools are installed |
| `itamae list` | List the available tools |
| `itamae logs` | View installation logs |
| `itamae version` | Show version information |

Every command has `--help`. Shell completions are available via `itamae completion <shell>`.

### install

Run the interactive installer:

```bash
itamae install
```

This launches a TUI (Terminal User Interface) where you can:

1. **Choose Category**: Select Core, Essentials, or Unverified packages
2. **Select Packages**: Pick which tools to install (Unverified only)
3. **Answer Questions**: Some tools need input, e.g. your Git name and email
4. **Review Plan**: Confirm your selections
5. **Monitor Progress**: Watch installation in real-time

#### Installation Categories

<details>
<summary><strong>Core</strong></summary>

Essential development tools installed automatically:
- git, GitHub CLI, kubectl, kubecolor, helm, yq, jq, task
- Node.js (from NodeSource, includes npm), python3, pipx, .NET SDK
- curl, wget, gnupg, ca-certificates, nala, fzf, fd, lsd, alacritty, WireGuard tools

</details>

<details>
<summary><strong>Essentials</strong></summary>

Common developer extras:
- stow, ripgrep, bat
- zoxide, starship, atuin
- rust, java (Temurin), maven, sdkman

</details>

<details>
<summary><strong>Unverified</strong></summary>

Additional tools with interactive multi-select:
- zsh, btop, ncdu, httpie, pass, tldr, zellij, chezmoi, bin
- Visual Studio Code, Ghostty (manual), Cascadia Code
- ansible, semgrep, ruby, meld, rofi, polybar, dunst, flameshot

</details>

#### Headless installs

Every prompt can be answered with a flag, so `itamae install` works from scripts, cloud-init and CI:

| Flag | Purpose |
|------|---------|
| `-c, --category` | Category to install: `core`, `essentials` or `unverified` |
| `-p, --plugins` | Comma-separated plugin IDs (see `itamae list --ids`) |
| `--all` | Install every plugin of the category (required for `unverified` without a terminal) |
| `-i, --input NAME=VALUE` | Answer a plugin question (repeatable). Environment variables with the same name work too |
| `-y, --yes` | Skip the confirmation prompt |
| `--plain` | Line-oriented output instead of the TUI (automatic when not attached to a terminal) |
| `--skip-verify` | Do not run each tool's `check()` after installing |

```bash
# Everything in core, no questions asked
GIT_USER_NAME="Jane Doe" GIT_USER_EMAIL="jane@example.com" itamae install --category core --yes

# A hand-picked set with plain output
itamae install -c unverified -p zsh,btop,ncdu --yes --plain
```

The command exits with status `1` if any tool failed to install or did not pass its post-install check, so you can chain it safely in scripts.

### check

Run every plugin's `check()` function and report what is installed:

```bash
itamae check                    # all categories
itamae check --category core    # one category
itamae check --plugins git,jq   # specific tools
```

The command exits non-zero when at least one selected tool is missing, which makes it handy for verifying a provisioning run.

### list

```bash
itamae list                      # every tool, grouped by category
itamae list --category essentials
itamae list -c unverified --ids  # plugin IDs only, one per line
```

### logs

View installation logs from previous runs:

```bash
# Show most recent log
itamae logs

# List all available logs
itamae logs --list

# View specific log file
itamae logs --file itamae-install-2024-01-15_14-30-45.log

# Follow log in real-time (like tail -f)
itamae logs --follow

# Filter log lines
itamae logs --grep "error"
itamae logs --grep "Phase"

# Clean up old logs
itamae logs --clean
```

Logs are saved to `/tmp/itamae-logs/` with timestamps and contain the complete output of every script and package manager command.

### version

```bash
itamae version
# or
itamae --version
```

## How It Works

### Installation Phases

1. **Repository Setup**: prerequisites (`curl`, `gnupg`, `ca-certificates`) are installed if missing, custom APT repositories (GitHub CLI, NodeSource, Adoptium, Microsoft) are configured, and package lists are refreshed once
2. **Batch Installation**: all APT packages are installed in one `nala` (or `apt-get`) command. If the batch fails, every package is retried individually so one unavailable package cannot fail the rest
3. **Individual Installation**: binary and manual plugins run one by one
4. **Verification**: each installed tool must pass its `check()` function before it counts as installed

### Performance Optimization

Itamae optimizes installation by:
- Batching all APT packages into a single command
- Using `nala` when available for parallel downloads
- Running repository setup before batch installation
- Keeping your `sudo` session alive during long installs so you are never asked for a password half way through

### Terminal User Interface

The TUI displays:
- **Left pane**: Package checklist with status icons and the current phase of each package
- **Right pane**: Scrollable installation log
- **Bottom pane**: Recent errors while installing, and the summary once complete

**Keyboard Navigation:**
- `↑/↓` or `j/k`: Scroll logs
- `PgUp/PgDown`: Page through logs
- `q`, `Enter` or `Ctrl+C`: Exit (after completion); `Ctrl+C` also aborts a running installation

The TUI uses the **Tokyo Night** color scheme.

## Troubleshooting

If you encounter issues:

1. Check the logs:
   ```bash
   itamae logs
   ```

2. Look for errors:
   ```bash
   itamae logs --grep ERROR
   ```

3. See what is and isn't installed:
   ```bash
   itamae check
   ```

For more help, see the [Troubleshooting Guide](/troubleshooting).
