---
title: Getting Started
weight: 2
---

# Getting Started

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

## Usage

To check which version of Itamae you have installed:
```bash
itamae version
# or
itamae --version
```

To run the installation wizard:
```bash
itamae install
```

## Available Plugins

The tool comes with various plugins. You can see the list during the installation process.

## How It Works

Itamae uses a Go binary to orchestrate shell scripts located in the `scripts/` directory. Each script is a plugin that knows how to install and remove a specific piece of software. The Go binary embeds these scripts and provides an interactive form interface for selecting optional plugins.

### Installation Process

1. **Interactive Selection:** Choose which additional tools to install using an interactive multi-select form
2. **Installation Plan:** Review the complete installation plan organized by method (APT, binary, manual)
3. **Batch Installation:** All APT packages are installed in a single optimized command using `nala` (or `apt-get`)
4. **Individual Installation:** Binary and manual installations run individually with live progress
