---
title: Itamae
archetype: home
---

# Itamae

{{% notice style="info" title="About Itamae" %}}
**Itamae** is a command-line tool written in Go that sets up a developer's Linux workstation using a plugin-based architecture.
{{% /notice %}}

## Features

- 🚀 **Fast Installation** - Optimized batch installation using `nala` or `apt-get` for parallel package downloads
- 🎨 **Beautiful TUI** - Modern terminal interface with Tokyo Night theme and real-time progress
- 🔌 **Plugin-Based** - Extensible architecture - easily add new tools via shell scripts
- 📦 **Smart Batching** - Single-command installation for all APT packages improves speed and dependency resolution
- 🐛 **Debug Logging** - Comprehensive logs saved to `/tmp/itamae-logs/` for troubleshooting
- 🎯 **Category System** - Organize tools into Core, Essentials, and Unverified categories

## Quick Start

Install Itamae with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

Then run:

```bash
itamae install
```

## What's Next?

### 📚 Installation

Learn different ways to install Itamae on your system.

[Get Started →]({{< relref "/docs/installation" >}})

### 🚀 Usage

Discover all the commands and features available.

[View Commands →]({{< relref "/docs/usage" >}})

### 👨‍💻 Development

Want to contribute? Learn how to add new plugins.

[Developer Guide →]({{< relref "/docs/developers" >}})
