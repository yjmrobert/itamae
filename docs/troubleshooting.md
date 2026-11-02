# Troubleshooting

Common issues and their solutions.

## Installation Issues

### Permission Denied

Do **not** run `itamae` with `sudo`: it would install user-level tools into root's home directory. Itamae asks for your password once up front (`sudo -v`) and keeps the session alive while it installs. If you see permission errors, make sure your user is allowed to use `sudo`:

```bash
sudo -v && echo "sudo works"
```

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
2. Inspect the complete command output in the debug log:
   ```bash
   itamae logs --grep "download\|fetch\|curl"
   ```
3. Re-run only the tools that failed:
   ```bash
   itamae check --category core          # shows what is missing
   itamae install -c core -p helm,kubectl --yes
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
