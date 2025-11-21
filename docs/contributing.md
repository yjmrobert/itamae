# Contributing to Itamae

Thank you for your interest in contributing to Itamae! This guide will help you get started.

## Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/itamae.git
   cd itamae
   ```
3. **Create a branch**:
   ```bash
   git checkout -b feature/my-new-feature
   ```
4. **Make your changes** and test them
5. **Commit** using conventional commits:
   ```bash
   git commit -m "feat: add awesome feature"
   ```
6. **Push** to your fork:
   ```bash
   git push origin feature/my-new-feature
   ```
7. **Open a Pull Request** on GitHub

## Code Guidelines

### Go Code Style

Follow standard Go conventions:

```bash
# Format your code
go fmt ./...

# Run static analysis
go vet ./...

# Run linter (if installed)
golangci-lint run
```

**Best Practices:**
- Use meaningful variable names
- Write comments for exported functions
- Keep functions small and focused
- Handle errors explicitly
- Use `go fmt` before committing

### Shell Script Style

For plugin scripts:

**Best Practices:**
- Use `#!/bin/bash` shebang
- Include complete metadata block
- Prefer `nala` over `apt-get` when available
- Echo progress messages
- Use the router pattern
- Create symlinks in `$HOME/.local/bin`

**Example:**
```bash
#!/bin/bash
#
# METADATA
# NAME: My Tool
# DESCRIPTION: A helpful tool
# INSTALL_METHOD: apt
# PACKAGE_NAME: my-tool
#

install() {
    echo "Installing My Tool..."
    if command -v nala &> /dev/null; then
        sudo nala install -y my-tool
    else
        sudo apt-get install -y my-tool
    fi
    echo "✅ My Tool installed."
}

remove() {
    echo "Removing My Tool..."
    sudo apt-get purge -y my-tool
    echo "✅ My Tool removed."
}

check() {
    command -v my-tool &> /dev/null
}

# --- ROUTER ---
case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
```

## Commit Convention

**IMPORTANT**: Itamae uses [Conventional Commits](https://www.conventionalcommits.org/) for automated versioning.

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description | Version Impact |
|------|-------------|----------------|
| `feat:` | New feature | Minor bump (0.x.0) |
| `fix:` | Bug fix | Patch bump (0.0.x) |
| `docs:` | Documentation only | Patch bump (0.0.x) |
| `style:` | Code style changes | Patch bump (0.0.x) |
| `refactor:` | Code refactoring | Patch bump (0.0.x) |
| `test:` | Adding tests | Patch bump (0.0.x) |
| `chore:` | Maintenance tasks | Patch bump (0.0.x) |
| `feat!:` or `BREAKING CHANGE:` | Breaking change | Major bump (x.0.0) |

### Examples

**Feature:**
```bash
git commit -m "feat: add support for zsh plugin"
```

**Bug Fix:**
```bash
git commit -m "fix: resolve symlink creation issue"
```

**Documentation:**
```bash
git commit -m "docs: update installation instructions"
```

**Breaking Change:**
```bash
git commit -m "feat!: redesign plugin metadata format

BREAKING CHANGE: plugins now require VERSION field in metadata"
```

## Testing

Before submitting a PR:

### Run Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v ./itamae -run TestGitPlugin
```

### Manual Testing

Test your changes in a clean environment:

```bash
# Build
./build.sh

# Test locally
./bin/itamae install
```

Consider testing in a VM or container:

```bash
# Using Docker
docker run -it ubuntu:22.04 bash

# Using Multipass
multipass launch -n test-itamae
multipass shell test-itamae
```

## Documentation

When making changes, update the documentation:

### Where to Update

- **Code changes**: Update inline comments
- **New features**: Update `docs/usage.md`
- **Plugin changes**: Update `docs/developers/adding-plugins.md`
- **Installation changes**: Update `docs/installation.md`

### Documentation Format

Documentation uses Markdown:

```markdown
# Page Title

Content goes here...

## Section

More content...
```

### Building Documentation

```bash
cd docs
hugo server
# Visit http://localhost:1313
```

## Pull Request Process

### Before Submitting

- [ ] Code follows style guidelines
- [ ] Tests pass: `go test ./...`
- [ ] Code is formatted: `go fmt ./...`
- [ ] Documentation is updated
- [ ] Commits follow conventional format
- [ ] Branch is up to date with master

### PR Template

When creating a PR, include:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (fix)
- [ ] New feature (feat)
- [ ] Breaking change (feat!)
- [ ] Documentation (docs)

## Testing
Describe how you tested the changes

## Checklist
- [ ] Tests pass
- [ ] Code formatted
- [ ] Documentation updated
- [ ] Conventional commits used
```

### Review Process

1. **Automated Checks**: CI runs tests and linting
2. **Code Review**: Maintainer reviews code
3. **Feedback**: Address any requested changes
4. **Approval**: PR is approved
5. **Merge**: Maintainer merges PR

## Adding New Plugins

See the [Adding Plugins Guide](/developers/adding-plugins) for detailed instructions.

**Quick checklist:**
- [ ] Script in appropriate directory (`core/`, `essentials/`, or `unverified/`)
- [ ] Complete metadata block
- [ ] `install()`, `remove()`, and `check()` functions
- [ ] Router case statement
- [ ] Test added to `itamae/main_test.go`
- [ ] Documentation updated

## Getting Help

Need help contributing?

- **GitHub Discussions**: Ask questions and discuss ideas
- **GitHub Issues**: Report bugs or request features
- **Documentation**: Check existing docs
- **Examples**: Look at existing plugins for patterns

## Code of Conduct

Be respectful and professional:

- Use welcoming and inclusive language
- Be respectful of differing viewpoints
- Accept constructive criticism gracefully
- Focus on what's best for the community

## Recognition

Contributors are recognized in:
- `CONTRIBUTORS.md` file
- GitHub contributors page
- Release notes

Thank you for contributing to Itamae! 🎉
