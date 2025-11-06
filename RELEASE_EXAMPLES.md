# Release Script Usage Examples

## Basic Usage

Simply run the script from the repository root:

```bash
./release.sh
```

The script will guide you through the process interactively.

## Example Output

### Scenario: First Release (v0.1.0)

Starting with no tags and several `feat:` commits:

```
========================================
Itamae Release Script
========================================
ℹ Running pre-flight checks...
✓ Working tree is clean

========================================
Running Tests
========================================
ℹ Formatting code...
✓ Code formatted
ℹ Running vet...
✓ Vet passed
ℹ Running tests...
✓ All tests passed

========================================
Version Analysis
========================================
ℹ Latest tag: v0.0.0
ℹ Current version: 0.0.0
ℹ Analyzing commits since v0.0.0...

  - feat: Enhance installation process and add version command support
  - feat(install): run itamae install even if already installed
  - feat: Update dependencies and enhance installation process
  - feat: Enhance plugin metadata and installation process
  - feat: Replace bootstrap script with install script
  - docs: Update documentation
  - fix: resolve installation error

ℹ Commit summary:
  Breaking changes: 0
  Features:         5
  Fixes:            1
  Other:            1

✓ Determined bump type: minor
✓ New version: v0.1.0

⚠ This will create and push tag v0.1.0
Continue with release? (y/N): y

========================================
Building Release
========================================
ℹ Building v0.1.0...
Building itamae v0.1.0 (e078d09)...
✅ Build complete: bin/itamae
   Version: v0.1.0
   Commit:  e078d09
   Date:    2025-11-06T16:00:00Z

ℹ Testing built binary...
itamae version v0.1.0
  commit: e078d09
  built:  2025-11-06T16:00:00Z
✓ Binary built and tested successfully

========================================
Creating Release
========================================
ℹ Creating git tag...
✓ Tag created: v0.1.0
ℹ Pushing to origin...
✓ Pushed branch to origin
ℹ Pushing tag to origin...
✓ Pushed tag to origin

========================================
Release Complete!
========================================

✓ Released v0.1.0

ℹ GitHub Actions will now build the release artifacts.
ℹ View the release at: https://github.com/yjmrobert/itamae/releases/tag/v0.1.0
```

### Scenario: Patch Release (v0.1.1)

With only bug fixes since v0.1.0:

```
ℹ Latest tag: v0.1.0
ℹ Current version: 0.1.0
ℹ Analyzing commits since v0.1.0...

  - fix: resolve zsh configuration issue
  - fix: correct install script path handling
  - docs: update README

ℹ Commit summary:
  Breaking changes: 0
  Features:         0
  Fixes:            2
  Other:            1

✓ Determined bump type: patch
✓ New version: v0.1.1
```

### Scenario: Minor Release (v0.2.0)

With new features since v0.1.1:

```
ℹ Latest tag: v0.1.1
ℹ Current version: 0.1.1
ℹ Analyzing commits since v0.1.1...

  - feat: add Neovim plugin
  - feat: add Tmux plugin
  - fix: improve error handling

ℹ Commit summary:
  Breaking changes: 0
  Features:         2
  Fixes:            1
  Other:            0

✓ Determined bump type: minor
✓ New version: v0.2.0
```

### Scenario: Major Release (v1.0.0)

With breaking changes:

```
ℹ Latest tag: v0.2.0
ℹ Current version: 0.2.0
ℹ Analyzing commits since v0.2.0...

  - feat!: change plugin API to use new interface
  - feat: add plugin versioning support
  - fix: resolve compatibility issues

ℹ Commit summary:
  Breaking changes: 1
  Features:         1
  Fixes:            1
  Other:            0

✓ Determined bump type: major
✓ New version: v1.0.0
```

## Pre-Release Checklist

Before running `./release.sh`:

1. ✅ All changes committed
2. ✅ Working tree clean
3. ✅ Tests passing locally
4. ✅ On master/main branch
5. ✅ Commits follow conventional format

## What Happens After Release

Once you push the tag:

1. **GitHub Actions Triggered**: The workflow in `.github/workflows/release.yml` starts
2. **Multi-Platform Builds**: Binaries are built for:
   - Linux AMD64
   - Linux ARM64
   - macOS AMD64
   - macOS ARM64
3. **Checksums Generated**: SHA256 checksums for all binaries
4. **GitHub Release Created**: Automatic release with:
   - All binary artifacts
   - Checksums file
   - Auto-generated release notes
   - Installation instructions

## Canceling a Release

If you need to cancel after creating the tag but before pushing:

```bash
# Delete local tag
git tag -d v0.1.0

# Don't push to origin
```

If you already pushed:

```bash
# Delete local tag
git tag -d v0.1.0

# Delete remote tag
git push origin :refs/tags/v0.1.0

# Delete the GitHub release manually from the web interface
```

## Troubleshooting

### "Not on master/main branch"

Switch to the main branch:
```bash
git checkout master
```

### "You have uncommitted changes"

Commit or stash your changes:
```bash
git status
git add .
git commit -m "feat: your commit message"
# or
git stash
```

### "No commits found since last tag"

You need at least one commit to create a release:
```bash
git log HEAD...<last-tag>
```

### Tests Failing

Fix the issues before releasing:
```bash
go fmt ./...
go vet ./...
go test ./...
```

## Manual Override

If you need to create a specific version (not recommended):

```bash
# Build with specific version
VERSION=v1.0.0 ./build.sh

# Create tag manually
git tag -a v1.0.0 -m "Release v1.0.0"

# Push
git push origin master
git push origin v1.0.0
```
