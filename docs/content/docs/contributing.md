---
title: Contributing
weight: 3
---

# Contributing to Itamae

Contributions are welcome!

If you would like to add a new tool to Itamae, please see the [Developer Guide](/docs/developers/) for detailed instructions.

## Linting and Testing

Before submitting a pull request, please run the linter and tests to ensure that your changes are correct and that you have not introduced any regressions.

```bash
go fmt ./...
go vet ./...
go test ./...
```

## Commit Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated version management:

- `feat:` - New feature (minor version bump)
- `fix:` - Bug fix (patch version bump)
- `feat!:` or `BREAKING CHANGE:` - Breaking change (major version bump)
- `docs:`, `chore:`, `style:`, `refactor:`, `test:` - Other changes (patch version bump)

Example:
```bash
git commit -m "feat: add support for Neovim plugin"
git commit -m "fix: resolve path issue in zsh script"
```
