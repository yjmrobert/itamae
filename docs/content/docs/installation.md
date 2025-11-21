---
title: Installation
weight: 10
---

# Installation

There are several ways to install Itamae on your Linux system.

## Quick Install (Recommended)

{{% notice style="info" title="Quick Install" %}}
**One-line installation** - Copy and paste this command to install Itamae:
{{% /notice %}}

```bash
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

This will:
- ✅ Download the latest release
- ✅ Install it to `/usr/local/bin`
- ✅ Make it available system-wide

{{% notice style="warning" title="Security Note" %}}
Always review scripts before piping to bash. You can inspect the install script at: [install.sh](https://github.com/yjmrobert/itamae/blob/master/install.sh)
{{% /notice %}}

## Install Specific Version

You can install a specific version by setting the `ITAMAE_VERSION` environment variable:

```bash
ITAMAE_VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

## Install with Go

If you have Go 1.24+ installed:

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
ln -sf ~/source/repos/itamae/bin/itamae ~/bin/itamae

# Add ~/bin to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/bin:$PATH"

# Verify
which itamae
itamae version
```

This way, each time you run `./build.sh`, the symlink automatically points to the latest build.

## Verify Installation

Check that Itamae is installed correctly:

```bash
itamae version
```

## Next Steps

Now that you have Itamae installed, learn how to use it:

[Usage Guide →]({{% relref "usage" %}})
