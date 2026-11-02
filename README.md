# Itamae

Itamae is a command-line tool written in Go that sets up a developer's Linux (Debian/Ubuntu) workstation. It uses a plugin-based architecture to install, verify and remove software, and is built on the [Charm](https://charm.sh/) ecosystem (Bubble Tea, Bubbles, Lip Gloss, Huh, Log and Fang) with [Cobra](https://github.com/spf13/cobra) for the CLI.

Full documentation lives in the [`docs/`](docs/) directory and is published as a VitePress site.

## Installation

### Quick Install (Latest Version)

```bash
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

The script installs prerequisites, downloads the latest release to `/usr/local/bin` (upgrading an existing copy when needed) and starts the interactive installer.

### Unattended bootstrap

Everything the installer asks can be answered up front, which makes a one-liner usable from provisioning scripts:

```bash
ITAMAE_INSTALL_ARGS="--category core --yes" \
GIT_USER_NAME="Jane Doe" GIT_USER_EMAIL="jane@example.com" \
  curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

| Variable | Purpose |
|----------|---------|
| `ITAMAE_VERSION` | Install a specific release (e.g. `v1.0.0`) instead of the latest |
| `ITAMAE_INSTALL_DIR` | Where to put the binary (default `/usr/local/bin`) |
| `ITAMAE_INSTALL_ARGS` | Arguments for the `itamae install` run that follows |
| `ITAMAE_SKIP_RUN=1` | Only install the binary, do not run `itamae install` |
| `ITAMAE_SKIP_PREREQS=1` | Skip `apt-get update` and the git/nala installation |

### Install with Go

```bash
go install github.com/yjmrobert/itamae@latest
```

### Build from Source

```bash
git clone https://github.com/yjmrobert/itamae.git
cd itamae
./build.sh
sudo mv bin/itamae /usr/local/bin/
```

## Commands

| Command | Description |
|---------|-------------|
| `itamae install` | Install a category of tools (interactive TUI, or headless with flags) |
| `itamae check` | Run every plugin's `check()` and report what is installed; exits non-zero if anything is missing |
| `itamae list` | List the available tools, grouped by category (`--ids` prints plugin IDs only) |
| `itamae logs` | View, follow, filter or clean the installation logs in `/tmp/itamae-logs/` |
| `itamae version` / `--version` | Show version information |

### Headless installs

```bash
# Install every core tool without prompts
itamae install --category core --yes

# Pick specific tools and answer plugin questions on the command line
itamae install -c unverified --plugins zsh,btop --yes --plain \
  --input GIT_USER_NAME="Jane Doe"

# Verify afterwards
itamae check --category core
```

`--plain` switches to line-oriented output, and is selected automatically when the output is not a terminal. Values for plugin questions (such as `GIT_USER_NAME`) are also read from environment variables of the same name. The command exits with a non-zero status when any tool fails to install or does not pass its post-install check.

## How It Works

Each tool is a shell script embedded into the binary with a small metadata header (`NAME`, `INSTALL_METHOD`, `PACKAGE_NAME`, ...). Scripts live in three categories:

- **core** - essential tools, installed together
- **essentials** - common developer extras, installed together
- **unverified** - optional tools you pick individually

An installation runs in phases:

1. **Repository setup** - prerequisites (`curl`, `gnupg`, `ca-certificates`) are installed if missing, custom APT repositories are configured, then package lists are refreshed once.
2. **Batch APT install** - every APT package is installed in a single `nala` (or `apt-get`) command. If the batch fails, packages are retried one at a time so a single unavailable package does not sink the rest.
3. **Individual installs** - binary and manual plugins run their own `install()` functions.
4. **Verification** - every installed tool's `check()` function must pass before it is reported as successful.

Script output is captured in a timestamped debug log (`itamae logs`), and the TUI shows a live checklist, log pane and error pane in the Tokyo Night theme.

## Continuous Integration

- **CI** (`.github/workflows/ci.yml`) - gofmt, `go vet`, golangci-lint, unit tests with the race detector, ShellCheck over every script, and a docs build.
- **Integration** (`.github/workflows/integration.yml`) - for every supported distribution (Ubuntu 22.04/24.04, Debian 12/13), both package managers (`nala`, `apt-get`) and every category, a fresh container runs `itamae install` headlessly, verifies each tool with `itamae check`, confirms the expected package manager was used, and removes everything again through the plugin scripts. A second job exercises every code path of `install.sh` against a local release mirror.
- **Release** (`.github/workflows/release.yml`) - builds and publishes binaries when a `v*` tag is pushed.
- Dependabot keeps Go modules, GitHub Actions and the docs dependencies current.

## Contributing

Contributions are welcome! See the [Developer Guide](DEVELOPERS.md) and the [docs](docs/) for how to add a tool. Commits follow [Conventional Commits](https://www.conventionalcommits.org/); releases are cut with `./release.sh` (see [RELEASE.md](RELEASE.md)).
