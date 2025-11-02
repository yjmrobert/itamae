# Contributing to Itamae

Contributions are welcome!

If you would like to add a new tool to Itamae, please see the [Developer Guide](DEVELOPERS.md) for detailed instructions.

## Linting and Testing

Before submitting a pull request, please run the linter and tests to ensure that your changes are correct and that you have not introduced any regressions.

```bash
go fmt ./...
go vet ./...
go test ./...
```
