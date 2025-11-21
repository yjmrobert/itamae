---
title: Usage Guide
weight: 20
---

# Usage Guide

Learn how to use Itamae to set up your development environment.

## Commands

### install

Run the interactive installer:

```bash
itamae install
```

This launches a TUI (Terminal User Interface) where you can:

1. **Choose Category**: Select Core, Essentials, or Unverified packages
2. **Select Packages**: Pick which tools to install
3. **Review Plan**: Confirm your selections
4. **Monitor Progress**: Watch installation in real-time

#### Installation Categories

{{% expand title="Core" %}}
Essential development tools installed automatically:
- git, kubectl, helm
- nodejs, npm
- python3, pipx
- And more...
{{% /expand %}}

{{% expand title="Essentials" %}}
Common developer extras:
- stow, ripgrep, bat
- zoxide, starship, atuin
- rust, java, maven, sdkman
{{% /expand %}}

{{% expand title="Unverified" %}}
Additional tools with interactive multi-select:
- Choose specific packages
- Customizable selection
{{% /expand %}}

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
itamae logs --grep "Phase 1"

# Clean up old logs
itamae logs --clean
```

Logs are saved to `/tmp/itamae-logs/` with timestamps.

### version

Display version information:

```bash
itamae version
# or
itamae --version
```

## How It Works

### Installation Process

#### Phase Breakdown

1. **Category Selection**: Choose between Core, Essentials, or Unverified
2. **Package Selection**: Pick specific tools (Unverified only)
3. **Confirmation**: Review and confirm
4. **Repository Setup** (Phase 0): Add custom repositories, run single `apt-get update`
5. **Batch Installation** (Phase 1): Install all APT packages in one optimized command
6. **Individual Installation** (Phase 2): Install binary/manual packages one by one

### Performance Optimization

Itamae optimizes installation by:
- Batching all APT packages into a single command
- Using `nala` when available for parallel downloads
- Running repository setup before batch installation
- Significantly faster than installing packages one-by-one

### Terminal User Interface

The TUI displays:
- **Left pane**: Package checklist with status icons
- **Right pane**: Scrollable installation logs
- **Bottom pane**: Error messages (when failures occur)

**Keyboard Navigation:**
- `↑/↓` or `j/k`: Scroll logs
- `PgUp/PgDown`: Page through logs
- `q` or `Ctrl+C`: Exit (after completion)

The TUI uses the **Tokyo Night** color scheme for a modern, readable appearance.

## Troubleshooting

If you encounter issues:

1. Check the logs:
   ```bash
   itamae logs
   ```

2. Look for errors:
   ```bash
   itamae logs --grep error
   ```

3. Follow the most recent installation:
   ```bash
   itamae logs --follow
   ```

For more help, see the [Troubleshooting Guide]({{< relref "troubleshooting" >}}).
