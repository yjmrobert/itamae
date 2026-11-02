# Contributing to Itamae

Contributions are welcome!

If you would like to add a new tool to Itamae, please see the [Developer Guide](DEVELOPERS.md) for the plugin contract and the testing workflow.

## Before opening a pull request

```bash
gofmt -l .                    # must print nothing
go vet ./...
go test -race ./...
golangci-lint run
shellcheck -S warning itamae/scripts/*/*.sh install.sh
```

The CI workflow runs the same checks, and the Integration workflow installs every plugin category for real on Ubuntu 22.04/24.04 and Debian 12/13 with both `nala` and `apt-get`.

Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/), e.g. `feat: add neovim plugin` or `fix(installer): handle missing curl`.
