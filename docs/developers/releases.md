# Release Process

Learn how to create releases for Itamae using the automated release system.

## Overview

Itamae uses a **fully automated release process** that:
- Analyzes conventional commits to determine version bump
- Runs tests and linting
- Builds binaries for multiple platforms
- Creates GitHub release with changelog
- Uploads release artifacts

## Prerequisites

Before creating a release, ensure:

1. **Conventional Commits**: All commits follow [Conventional Commits](https://www.conventionalcommits.org/) format
2. **Tests Pass**: All tests are passing
3. **Clean Working Tree**: No uncommitted changes
4. **GitHub Token**: `GITHUB_TOKEN` environment variable set (for GitHub CLI)

## Conventional Commits

Itamae uses conventional commits to automatically determine version bumps:

### Commit Types

| Type | Version Bump | Example |
|------|-------------|---------|
| `feat:` | **Minor** (0.x.0) | `feat: add neovim plugin` |
| `fix:` | **Patch** (0.0.x) | `fix: resolve path issue in installer` |
| `feat!:` or `BREAKING CHANGE:` | **Major** (x.0.0) | `feat!: redesign TUI interface` |
| `docs:`, `chore:`, `style:`, `refactor:`, `test:` | **Patch** (0.0.x) | `docs: update README` |

### Commit Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Examples

**Feature (minor bump):**
```bash
git commit -m "feat: add support for Ghostty terminal"
```

**Bug Fix (patch bump):**
```bash
git commit -m "fix: correct kubectl version detection"
```

**Breaking Change (major bump):**
```bash
git commit -m "feat!: change plugin metadata format

BREAKING CHANGE: plugins now require VERSION metadata field"
```

**Documentation (patch bump):**
```bash
git commit -m "docs: add TUI development guide"
```

## Creating a Release

### Automated Release (Recommended)

Use the release script:

```bash
./release.sh
```

This script will:

1. **Verify Prerequisites**
   - Check for uncommitted changes
   - Verify tests pass
   - Ensure on master branch

2. **Determine Version**
   - Analyze commits since last release
   - Calculate next version based on commit types
   - Display version bump rationale

3. **Update Files**
   - Update version in `main.go`
   - Create/update `CHANGELOG.md`
   - Commit version changes

4. **Create Release**
   - Tag commit with version
   - Push tag to GitHub
   - Trigger GitHub Actions build
   - Create GitHub release with changelog

5. **Publish Artifacts**
   - Build binaries for Linux (amd64, arm64)
   - Build binaries for macOS (amd64, arm64)
   - Build binaries for Windows (amd64)
   - Upload to GitHub release

### Manual Release

If you need to create a release manually:

```bash
# 1. Determine next version
VERSION="v1.2.3"

# 2. Update version in code
sed -i "s/version = \".*\"/version = \"$VERSION\"/" main.go

# 3. Run tests
go test ./...

# 4. Commit version bump
git add main.go
git commit -m "chore: bump version to $VERSION"

# 5. Create and push tag
git tag -a "$VERSION" -m "Release $VERSION"
git push origin master
git push origin "$VERSION"

# 6. Create GitHub release
gh release create "$VERSION" \
  --title "$VERSION" \
  --notes "Release notes here"
```

## Release Workflow

### GitHub Actions

The release process uses GitHub Actions (`.github/workflows/release.yml`):

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run tests
        run: go test ./...
      
      - name: Build binaries
        run: |
          # Linux amd64
          GOOS=linux GOARCH=amd64 go build -o itamae-linux-amd64
          
          # Linux arm64
          GOOS=linux GOARCH=arm64 go build -o itamae-linux-arm64
          
          # macOS amd64
          GOOS=darwin GOARCH=amd64 go build -o itamae-darwin-amd64
          
          # macOS arm64
          GOOS=darwin GOARCH=arm64 go build -o itamae-darwin-arm64
          
          # Windows amd64
          GOOS=windows GOARCH=amd64 go build -o itamae-windows-amd64.exe
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            itamae-linux-amd64
            itamae-linux-arm64
            itamae-darwin-amd64
            itamae-darwin-arm64
            itamae-windows-amd64.exe
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Version Management

### Version Format

Itamae uses [Semantic Versioning](https://semver.org/):

```
vMAJOR.MINOR.PATCH

Example: v1.2.3
```

- **MAJOR**: Breaking changes (incompatible API changes)
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Version in Code

Version is defined in `main.go`:

```go
var (
    version = "1.2.3"
    commit  = "unknown"
    date    = "unknown"
)
```

This is displayed by:

```bash
itamae version
# Output: Itamae v1.2.3 (commit: abc1234, built: 2024-11-21)
```

## Changelog

### Automatic Generation

The release script automatically generates `CHANGELOG.md`:

```markdown
# Changelog

## [1.2.3] - 2024-11-21

### Added
- Add support for Ghostty terminal
- Add neovim plugin

### Fixed
- Fix kubectl version detection
- Resolve path issue in installer

### Changed
- Update documentation
- Improve TUI layout
```

### Manual Editing

You can manually edit `CHANGELOG.md` before releasing:

```bash
# Edit changelog
vim CHANGELOG.md

# Commit changes
git add CHANGELOG.md
git commit -m "docs: update changelog"
```

## Release Checklist

Before releasing:

- [ ] All commits use conventional commit format
- [ ] Tests pass: `go test ./...`
- [ ] Linting passes: `go fmt ./...` and `go vet ./...`
- [ ] Documentation is up to date
- [ ] No uncommitted changes
- [ ] On master branch
- [ ] `GITHUB_TOKEN` is set

## Post-Release

After releasing:

1. **Verify Release**
   - Check GitHub releases page
   - Verify all artifacts are uploaded
   - Test installation from release

2. **Announce Release**
   - Update documentation
   - Post in discussions
   - Share on social media

3. **Monitor Issues**
   - Watch for bug reports
   - Respond to user feedback
   - Plan next release

## Hotfix Releases

For critical bug fixes:

```bash
# Create hotfix branch
git checkout -b hotfix/1.2.4

# Make fixes
git commit -m "fix: critical security issue"

# Merge to master
git checkout master
git merge hotfix/1.2.4

# Create release
./release.sh
```

## Troubleshooting

### Release Script Fails

**Problem**: Script reports uncommitted changes

**Solution**: Commit or stash changes
```bash
git status
git add .
git commit -m "chore: prepare for release"
```

**Problem**: Tests failing

**Solution**: Fix tests before releasing
```bash
go test -v ./...
```

**Problem**: GitHub token not set

**Solution**: Set `GITHUB_TOKEN`
```bash
export GITHUB_TOKEN="your_token_here"
```

### GitHub Actions Fails

**Problem**: Build fails on GitHub Actions

**Solution**: Check logs and test locally
```bash
# Test builds locally
GOOS=linux GOARCH=amd64 go build
GOOS=darwin GOARCH=arm64 go build
```

**Problem**: Release artifacts not uploading

**Solution**: Check GitHub token permissions and workflow file

## Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Go Cross Compilation](https://go.dev/doc/install/source#environment)

## Examples

See [RELEASE_EXAMPLES.md](https://github.com/yjmrobert/itamae/blob/master/RELEASE_EXAMPLES.md) for detailed examples of:
- Different commit types and their version impacts
- Multi-commit releases
- Breaking changes
- Hotfix scenarios
