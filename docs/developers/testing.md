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

Itamae uses a **mock-based testing approach** for shell script plugins. Due to environment limitations, full integration tests would require root access and could modify the system state. Instead, we use a controlled test environment.

### How It Works

The test suite (`itamae/main_test.go`) implements:

1. **Mock Environment Setup**: `TestMain()` function creates a temporary directory with mock commands
2. **Command Interception**: System commands like `sudo`, `apt-get`, `nala` are replaced with logging scripts
3. **Verification**: Tests verify that plugins call the correct commands by checking log files

### Test Structure

```go
func TestMain(m *testing.M) {
    // Setup mock environment
    tempDir := setupMockEnvironment()
    
    // Run tests
    code := m.Run()
    
    // Cleanup
    os.RemoveAll(tempDir)
    os.Exit(code)
}

func TestGitPlugin(t *testing.T) {
    // Execute plugin script
    output := executeScript("core/git.sh", "install")
    
    // Verify correct commands were called
    assertContains(t, output, "apt-get install -y git")
}
```

## Testing Your Plugin

When adding a new plugin, create a test that:

1. **Tests installation**: Verify `install` command works
2. **Tests removal**: Verify `remove` command works
3. **Tests detection**: Verify `check` command works
4. **Tests repository setup**: If applicable, verify `setup_repo` works

### Example Test

```go
func TestMyToolPlugin(t *testing.T) {
    script := "unverified/mytool.sh"
    
    t.Run("Install", func(t *testing.T) {
        output := executeScript(script, "install")
        if !strings.Contains(output, "apt-get install -y mytool") {
            t.Errorf("Expected install command, got: %s", output)
        }
    })
    
    t.Run("Remove", func(t *testing.T) {
        output := executeScript(script, "remove")
        if !strings.Contains(output, "apt-get purge -y mytool") {
            t.Errorf("Expected purge command, got: %s", output)
        }
    })
    
    t.Run("Check", func(t *testing.T) {
        // Mock successful check
        output := executeScript(script, "check")
        // Verify check logic
    })
}
```

## Manual Testing

For thorough testing, use a virtual machine or container:

### Using Docker

```bash
# Build test container
docker build -t itamae-test .

# Run installation
docker run -it itamae-test bash
./itamae install
```

### Using Vagrant

```bash
# Start Ubuntu VM
vagrant up

# SSH into VM
vagrant ssh

# Install and test Itamae
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
itamae install
```

### Using Multipass

```bash
# Launch Ubuntu instance
multipass launch -n itamae-test

# Shell into instance
multipass shell itamae-test

# Test installation
curl -fsSL https://raw.githubusercontent.com/yjmrobert/itamae/master/install.sh | bash
itamae install
```

## Debugging Tests

### Verbose Output

Run with verbose flag:

```bash
go test -v ./...
```

### Run Single Test

```bash
go test -v ./itamae -run TestGitPlugin
```

### Print Debug Info

Add print statements in tests:

```go
func TestDebug(t *testing.T) {
    output := executeScript("core/git.sh", "install")
    t.Logf("Output: %s", output)
}
```

## Continuous Integration

Tests run automatically on:
- Pull requests
- Commits to master branch
- Release builds

See `.github/workflows/` for CI configuration.

## Best Practices

1. **Test all plugin functions**: `install`, `remove`, `check`, and `setup_repo`
2. **Verify command accuracy**: Ensure exact commands are called
3. **Test error handling**: Verify plugins handle failures gracefully
4. **Mock external dependencies**: Don't rely on network or system state
5. **Keep tests fast**: Tests should complete quickly
6. **Test edge cases**: Empty inputs, missing commands, permission errors

## Linting

Run code quality checks:

```bash
# Format code
go fmt ./...

# Run static analysis
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run
```

## Pre-commit Checklist

Before committing:

- [ ] All tests pass: `go test ./...`
- [ ] Code is formatted: `go fmt ./...`
- [ ] No vet warnings: `go vet ./...`
- [ ] Plugin scripts have correct metadata
- [ ] Plugin scripts use router pattern
- [ ] New plugins have corresponding tests

## Next Steps

- [TUI Development](/developers/tui) - Working with the terminal UI
- [Release Process](/developers/releases) - Creating releases
