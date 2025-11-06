# Release Process

This document outlines how to create new releases for Itamae.

## Semantic Versioning

Itamae follows [Semantic Versioning](https://semver.org/) (SemVer):

- **Major version** (v2.0.0): Breaking changes to API or behavior
- **Minor version** (v1.1.0): New features or plugins added (backward compatible)
- **Patch version** (v1.0.1): Bug fixes (backward compatible)

## Conventional Commits

Itamae uses [Conventional Commits](https://www.conventionalcommits.org/) to automatically determine version bumps:

- `feat:` - New feature (triggers **minor** version bump)
- `fix:` - Bug fix (triggers **patch** version bump)
- `feat!:` or `BREAKING CHANGE:` - Breaking change (triggers **major** version bump)
- Other types (`docs:`, `chore:`, `style:`, `refactor:`, `test:`) - **patch** version bump

### Examples:

```bash
# Minor version bump (adds feature)
git commit -m "feat: add support for new plugin system"

# Patch version bump (fixes bug)
git commit -m "fix: resolve installation error on Ubuntu 24.04"

# Major version bump (breaking change)
git commit -m "feat!: change plugin configuration format"
# or
git commit -m "feat: change plugin API

BREAKING CHANGE: Plugin interface now requires version field"
```

## Creating a Release

### Automated Release (Recommended)

The easiest way to create a release is using the automated script:

```bash
./release.sh
```

This script will:
1. Check that you're on the master branch
2. Verify working tree is clean
3. Run tests (`go fmt`, `go vet`, `go test`)
4. Analyze commits since the last tag using conventional commit syntax
5. Automatically determine the version bump (major/minor/patch)
6. Build the new version
7. Create and push the git tag
8. Push changes to origin

### Manual Release

### Manual Release

If you prefer to create a release manually:

### 1. Prepare the Release

Ensure all changes are committed and tests pass:

```bash
go fmt ./...
go vet ./...
go test ./...
```

### 2. Create and Push a Git Tag

```bash
# Create an annotated tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# Push the tag to GitHub
git push origin v1.0.0
```

### 3. Automated Build

Once the tag is pushed, GitHub Actions will automatically:
- Build binaries for multiple platforms (Linux/macOS, AMD64/ARM64)
- Create checksums for verification
- Create a GitHub release with all artifacts
- Generate release notes from commits

### 4. Manual Build (Optional)

To build a release locally:

```bash
VERSION=v1.0.0 ./build.sh
```

This creates a versioned binary in `bin/itamae`.

## Installation Methods

### Using install.sh (Latest)

```bash
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

### Using install.sh (Specific Version)

```bash
ITAMAE_VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
```

### Using go install (Latest)

```bash
go install github.com/yjmrobert/itamae@latest
```

### Using go install (Specific Version)

```bash
go install github.com/yjmrobert/itamae@v1.0.0
```

### From GitHub Release

Download the appropriate binary from the [releases page](https://github.com/yjmrobert/itamae/releases):

```bash
# Example for Linux AMD64
curl -L -o itamae https://github.com/yjmrobert/itamae/releases/download/v1.0.0/itamae-linux-amd64
chmod +x itamae
sudo mv itamae /usr/local/bin/
```

## Version Information

Users can check their installed version:

```bash
itamae version
itamae --version
itamae -v
```

Output format:
```
itamae version v1.0.0
  commit: abc123f
  built:  2025-11-06T15:00:00Z
```

## Release Checklist

Before running the release script, ensure:

- [ ] All changes are committed
- [ ] Working tree is clean (`git status`)
- [ ] On master/main branch
- [ ] Commits follow conventional commit format
- [ ] Tests pass locally: `go test ./...`

Then simply run:
```bash
./release.sh
```

### What the Script Does

The `release.sh` script automates the entire release process:

**Pre-flight Checks:**
- Verifies you're on the master/main branch
- Ensures working tree is clean (no uncommitted changes)

**Code Quality:**
- Formats code with `go fmt ./...`
- Runs static analysis with `go vet ./...`
- Executes all tests with `go test ./...`

**Version Bumping:**
- Fetches the latest git tag
- Analyzes all commits since last tag
- Counts conventional commit types (breaking, feat, fix)
- Determines appropriate version bump (major/minor/patch)
- Calculates and displays the new version

**Build & Release:**
- Builds the binary with version information
- Tests the built binary
- Creates an annotated git tag
- Pushes branch and tag to origin
- Triggers GitHub Actions for automated release artifacts

**Interactive:**
- Shows commit summary before releasing
- Displays version bump calculation
- Requires confirmation before proceeding

### Manual Release Checklist

- [ ] Update CHANGELOG.md with release notes (if exists)
- [ ] Run tests: `go test ./...`
- [ ] Format code: `go fmt ./...`
- [ ] Check for issues: `go vet ./...`
- [ ] Commit all changes
- [ ] Create annotated git tag
- [ ] Push tag to GitHub
- [ ] Verify GitHub Actions build completes
- [ ] Test installation from release
- [ ] Announce release (if applicable)

## First Release

For the initial v1.0.0 release:

```bash
# Ensure you're on the main/master branch
git checkout master

# Create the first tag
git tag -a v1.0.0 -m "Initial stable release

Features:
- Interactive plugin selection with Huh forms
- Batch APT package installation
- Support for 25+ development tools
- Automated setup script
- Version management system
"

# Push to GitHub
git push origin v1.0.0
```

The GitHub Actions workflow will handle the rest!
