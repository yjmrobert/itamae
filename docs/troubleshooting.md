# Troubleshooting

Common issues and their solutions.

## Installation Issues

### Permission Denied

If you see permission errors:

```bash
sudo itamae install
```

Some packages require root privileges for installation.

### Package Not Found

If a package can't be found:

1. Update package lists:
   ```bash
   sudo apt-get update
   ```

2. Check if the repository was added:
   ```bash
   itamae logs --grep "repository"
   ```

### Network Errors

For network-related issues:

1. Check your internet connection
2. Try again with verbose logging:
   ```bash
   itamae install
   itamae logs --grep "download\|fetch\|curl"
   ```

## Debug Logging

:::info Automatic Logging
Itamae automatically saves detailed logs to `/tmp/itamae-logs/` for every installation run.
:::

### Viewing Logs

```bash
# Most recent log
itamae logs

# List all logs
itamae logs --list

# Search for errors
itamae logs --grep error
```

### Log Contents

Logs include:
- All commands executed
- Complete command output
- Phase transitions and timing
- Package status updates
- Error messages with context

### Log Format

```
[HH:MM:SS.mmm] Log message
```

The `itamae logs` command colorizes output:
- **Red**: Errors
- **Green**: Success messages
- **Blue**: Phase transitions
- **Gray**: Timestamps

## Common Problems

### Binary Not Found After Installation

Add the installation directory to your PATH:

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH="$HOME/.local/bin:$PATH"

# Reload shell
source ~/.bashrc
```

### Conflicting Packages

If packages conflict:

1. Remove the conflicting package:
   ```bash
   sudo apt-get remove <package-name>
   ```

2. Run itamae again:
   ```bash
   itamae install
   ```

## Getting Help

If you can't resolve an issue:

1. Check existing [GitHub Issues](https://github.com/yjmrobert/itamae/issues)
2. Create a new issue with:
   - Your OS version
   - Itamae version (`itamae version`)
   - Relevant log output (`itamae logs`)
   - Steps to reproduce

## Reporting Bugs

When reporting bugs, include:

```bash
# System information
uname -a
lsb_release -a

# Itamae version
itamae version

# Recent log
itamae logs
```
