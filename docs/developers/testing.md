# Testing

Learn how to test your Itamae plugins and contributions.

## Running Tests

Run all tests with:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run tests for a specific package:

```bash
go test ./itamae/...
go test ./cmd/...
```

## Test Coverage

Generate coverage report:

```bash
go test -cover ./...
```

Generate detailed coverage HTML report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Testing Approach

Itamae has three layers of tests.

### 1. Unit tests (`go test`)

- **Installer engine** (`itamae/engine_test.go`): fake plugin scripts and a mocked `sudo` drive the `Installer` through every branch: repository setup, prerequisites, `nala` vs `apt-get`, batch install, the per-package fallback when a batch fails, post-install hooks, verification failures, manual plugins, cancellation and the plain-text reporter.
- **Plugin scripts** (`itamae/main_test.go`): every embedded script runs its `install`, `remove`, `setup_repo` and `post_install` actions with `sudo`, `nala`, `apt-get`, `curl`, ... replaced by logging stubs. The test asserts the expected commands are issued, that unknown actions print usage, and that the catalogue of plugins (`expectedMetadata`) matches the scripts exactly, so adding or removing a script is always a deliberate change.
- **CLI** (`cmd/*_test.go`): flag handling, input resolution from `--input` and the environment, and the headless error paths.

```go
func TestInstallerBatchFailureFallsBackToIndividualInstalls(t *testing.T) {
    env := newEngineEnv(t, true)              // mocked sudo + nala on PATH
    plugins := []ToolPlugin{
        env.plugin("goodpkg", "apt", 0, 0, "post"),
        env.plugin("badpkg", "apt", 0, 0),    // the sudo mock fails for "badpkg"
    }
    sender, result := env.run(plugins, nil)
    // goodpkg succeeds individually, badpkg is reported as failed
}
```

### 2. install.sh end-to-end (`test/install-sh/test.sh`)

A local HTTP server stands in for GitHub releases so every path of `install.sh` is exercised deterministically: fresh install, already up to date, upgrade, downgrade with `ITAMAE_VERSION`, replacing a binary with an unknown version, a failed download that keeps the previous binary, the non-interactive guidance, `ITAMAE_INSTALL_ARGS` running a real headless install, unsupported architectures and the prerequisite step.

```bash
test/install-sh/test.sh
```

### 3. Container matrix (`test/integration/`)

The Integration workflow builds the binary once and then, for every combination of

- distribution: `ubuntu:22.04`, `ubuntu:24.04`, `debian:12`, `debian:13`
- package manager: `nala`, `apt-get`
- category: `core`, `essentials`, `unverified`

starts a fresh container with a regular user and passwordless sudo and runs `test/integration/run.sh`, which:

1. installs or removes `nala` to set up the permutation
2. runs `itamae install --category <c> --yes --plain` (with `--all` for unverified)
3. runs `itamae check` and requires every installable tool to pass
4. inspects the debug log to confirm the expected package manager was used
5. removes every tool through its script and asserts `check()` fails afterwards

Run one permutation locally with Docker:

```bash
docker build --build-arg BASE_IMAGE=debian:12 -t itamae-test -f test/integration/Dockerfile .
go build -o dist/itamae .
docker run --rm -e CATEGORY=essentials -e PKG_MANAGER=apt-get \
  -v "$PWD:/src:ro" -v "$PWD/dist/itamae:/usr/local/bin/itamae:ro" \
  itamae-test bash /src/test/integration/run.sh
```

## Testing Your Plugin

When adding a new plugin:

1. Add it to `expectedMetadata` and `pluginAssertions` in `itamae/main_test.go`; the suite fails for unregistered scripts.
2. Run `go test ./itamae/ -run 'TestInstall|TestUninstall|TestMetadataParsing'`.
3. Make sure `check()` fails after `remove()`; the container matrix asserts it.
4. Run the relevant container permutation locally, or push and let the Integration workflow do it.

## Debugging Tests

```bash
go test -v ./itamae -run TestInstall/install_git     # one script action
go test -v ./itamae -run TestInstaller               # engine tests
itamae logs                                          # after a real run
```

## Continuous Integration

- **CI** (`ci.yml`): gofmt, `go mod tidy` check, `go vet`, `go test -race`, golangci-lint, ShellCheck over every script, and a docs build. Runs on every push and pull request.
- **Integration** (`integration.yml`): the container matrix and the `install.sh` suite. Runs on every push and pull request, weekly on a schedule to catch upstream changes, and on demand.
- **Release** (`release.yml`): tests, then builds binaries for Linux and macOS when a `v*` tag is pushed.

## Linting

Run code quality checks:

```bash
# Format code
go fmt ./...

# Run static analysis
go vet ./...

# Run golangci-lint (configuration in .golangci.yml)
golangci-lint run

# Lint the shell scripts
shellcheck -S warning itamae/scripts/*/*.sh install.sh
```

## Pre-commit Checklist

Before committing:

- [ ] All tests pass: `go test ./...`
- [ ] Code is formatted: `go fmt ./...`
- [ ] No vet warnings: `go vet ./...`
- [ ] Lint is clean: `golangci-lint run` and `shellcheck -S warning itamae/scripts/*/*.sh`
- [ ] Plugin scripts have correct metadata and use the router pattern
- [ ] New plugins are registered in `itamae/main_test.go`

## Next Steps

- [TUI Development](/developers/tui) - Working with the terminal UI
- [Release Process](/developers/releases) - Creating releases
