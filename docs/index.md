---
layout: home

hero:
  name: "Itamae"
  text: "Linux Workstation Setup Tool"
  tagline: "A command-line tool written in Go that sets up a developer's Linux workstation using a plugin-based architecture"
  actions:
    - theme: brand
      text: Get Started
      link: /installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/yjmrobert/itamae

features:
  - icon: 🚀
    title: Fast Installation
    details: Optimized batch installation using nala or apt-get for parallel package downloads
  - icon: 🎨
    title: Beautiful TUI
    details: Modern terminal interface with cyber retro Tokyo Night theme and real-time progress
  - icon: 🔌
    title: Plugin-Based
    details: Extensible architecture - easily add new tools via shell scripts
  - icon: 📦
    title: Smart Batching
    details: Single-command installation for all APT packages improves speed and dependency resolution
  - icon: 🐛
    title: Debug Logging
    details: Comprehensive logs saved to /tmp/itamae-logs/ for troubleshooting
  - icon: 🎯
    title: Category System
    details: Organize tools into Core, Essentials, and Unverified categories
---

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

<div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem; margin-top: 2rem;">

<div>

### 📚 Installation

Learn different ways to install Itamae on your system.

[Get Started →](/installation)

</div>

<div>

### 🚀 Usage

Discover all the commands and features available.

[View Commands →](/usage)

</div>

<div>

### 👨‍💻 Development

Want to contribute? Learn how to add new plugins.

[Developer Guide →](/developers/)

</div>

</div>
