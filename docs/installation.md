# Installation

There are several ways to install Itamae on your Linux system.

## Quick Install (Recommended)

:::info Quick Install
**One-line installation** - Copy and paste this command to install Itamae:
:::

```bash
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

This will:
- ✅ Install the prerequisites (`git`, and `nala` where available)
- ✅ Download the latest release (or upgrade an existing installation)
- ✅ Install it to `/usr/local/bin`
- ✅ Start the interactive installer

:::warning Security Note
Always review scripts before piping to bash. You can inspect the install script at: [install.sh](https://github.com/yjmrobert/itamae/blob/master/install.sh)
:::

## Unattended Installation

The install script can be driven entirely by environment variables, which makes it suitable for provisioning scripts, cloud-init or dotfile bootstrappers:

```bash
ITAMAE_INSTALL_ARGS="--category core --yes" \
GIT_USER_NAME="Jane Doe" GIT_USER_EMAIL="jane@example.com" \
  curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

| Variable | Purpose |
|----------|---------|
| `ITAMAE_VERSION` | Install a specific release (e.g. `v1.0.0`) instead of the latest |
| `ITAMAE_INSTALL_DIR` | Where to put the binary (default `/usr/local/bin`) |
| `ITAMAE_INSTALL_ARGS` | Arguments for the `itamae install` run that follows the installation |
| `ITAMAE_SKIP_RUN=1` | Only install the binary and do not run `itamae install` |
| `ITAMAE_SKIP_PREREQS=1` | Skip `apt-get update` and the git/nala installation |

The script compares the installed version with the target and only downloads when an upgrade (or an explicit downgrade via `ITAMAE_VERSION`) is needed.

## Install Specific Version

```bash
ITAMAE_VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

## Install with Go

If you have Go 1.26+ installed:

```bash
# Latest version
go install github.com/yjmrobert/itamae@latest

# Specific version
go install github.com/yjmrobert/itamae@v1.0.0
```

The binary will be installed in `$HOME/go/bin`. Make sure this directory is in your `PATH`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

## Build from Source

For development or customization:

```bash
git clone https://github.com/yjmrobert/itamae.git
cd itamae
./build.sh
sudo mv bin/itamae /usr/local/bin/
```

### Development Setup

For local development, create a symlink instead:

```bash
# Build the project
./build.sh

# Create symlink in your local bin directory
mkdir -p ~/bin
ln -sf "$PWD/bin/itamae" ~/bin/itamae

# Add ~/bin to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/bin:$PATH"

# Verify
which itamae
itamae version
```

This way, each time you run `./build.sh`, the symlink automatically points to the latest build.

## Verify Installation

```bash
itamae version
itamae list
```

## Next Steps

Now that you have Itamae installed, learn how to use it:

[Usage Guide →](/usage)
