# Developer Guide

This is a condensed guide for contributors. The full, maintained version lives in the documentation site under [`docs/developers/`](docs/developers/).

## Architecture

```
main.go                 Entry point (cmd.Execute)
cmd/                    Cobra commands, run through Charm Fang
  root.go               Root command + fang.Execute (styled help, --version, man pages, completions)
  install.go            itamae install (interactive or headless)
  check.go              itamae check
  list.go               itamae list
  logs.go               itamae logs
  version.go            itamae version
itamae/
  plugin.go             Plugin metadata parsing and loading of the embedded scripts
  engine.go             Installer: phases, batching, fallback, verification (UI independent)
  exec.go               Running scripts and package manager commands
  run.go                TUI / plain-text front ends, PlainSender, PrintSummary
  forms.go              Huh forms (category, multi-select, confirm, inputs)
  tui_model.go/.._view.go  Bubble Tea v2 model and rendering
  theme.go              Tokyo Night palette (Lip Gloss v2)
  debug_logger.go       Timestamped debug log in /tmp/itamae-logs
  scripts/{core,essentials,unverified}/*.sh   The plugins
test/
  integration/          Dockerfile + run.sh used by the Integration workflow
  install-sh/           End-to-end tests for install.sh
```

The `Installer` in `engine.go` reports progress through the `Sender` interface. `*tea.Program` implements it for the TUI; `PlainSender` prints lines for `--plain`, CI and logs. Unit tests drive the installer with fake scripts and a mocked `sudo`.

## Plugin contract

A plugin is a bash script with a metadata header and four functions:

```bash
#!/bin/bash
#
# METADATA
# NAME: My Tool
# DESCRIPTION: One line about the tool.
# INSTALL_METHOD: apt            # apt | binary | manual
# PACKAGE_NAME: my-tool          # apt only
# REPO_SETUP: setup_repo         # optional, apt only
# POST_INSTALL: post_install     # optional, apt only
# REQUIRES: VAR|Prompt|default-command   # optional, repeatable
#

install() { ... }        # install the tool (used for binary/manual, and as the apt fallback)
remove()  { ... }        # remove it again
check()   { ... }        # exit 0 when the tool is usable

case "$1" in
    install) install ;;
    remove) remove ;;
    check) check ;;
    *) echo "Usage: $0 {install|remove|check}" && exit 1 ;;
esac
```

Rules of thumb:

- Prefer `nala` and fall back to `apt-get` for APT installs.
- Run package managers as `sudo env DEBIAN_FRONTEND=noninteractive ...`; `sudo` drops the variable otherwise and debconf blocks on prompts in non-interactive runs.
- Use `curl`, not `wget`, in `setup_repo()`: `curl`, `gnupg` and `ca-certificates` are guaranteed to be present before repository setup runs.
- Detect the CPU architecture (`uname -m`) when downloading binaries; releases ship for amd64 and arm64.
- Install user-level binaries into `~/.local/bin`; it is on `PATH` when scripts and `check()` run.
- Install any extra tools your script needs (for example `unzip`) inside the script.
- `check()` must fail after `remove()`; the integration tests assert this.

After adding a script, register it in `itamae/main_test.go` (`expectedMetadata` and `pluginAssertions`). The tests run every script against mocked commands and fail for unknown plugins.

## Development loop

```bash
go build ./... && go vet ./...
go test -race ./...
golangci-lint run
shellcheck -S warning itamae/scripts/*/*.sh install.sh
go run . install --category core --yes --plain      # headless run (needs sudo)
test/install-sh/test.sh                             # install.sh end-to-end
```

The container matrix (`test/integration/run.sh`) can be run locally with Docker:

```bash
docker build --build-arg BASE_IMAGE=ubuntu:24.04 -t itamae-test -f test/integration/Dockerfile .
go build -o dist/itamae .
docker run --rm -e CATEGORY=core -e PKG_MANAGER=nala \
  -v "$PWD:/src:ro" -v "$PWD/dist/itamae:/usr/local/bin/itamae:ro" \
  itamae-test bash /src/test/integration/run.sh
```

## Commits and releases

Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `docs:`, `chore:`, ...). `./release.sh` derives the next version from the commits and pushes a tag; the Release workflow builds the binaries.
