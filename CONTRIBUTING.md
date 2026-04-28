# Contributing to SpiceDB Go Library

Thank you for your interest in contributing! We welcome contributions of all kinds.

## Getting Started

### Prerequisites

- Go 1.20 or later
- Docker (for running SpiceDB locally)
- Git

### Setting Up Development Environment

1. Clone the repository:
```bash
git clone https://github.com/Yasiruofficial/events-authz.git
cd spicedb-go
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests to ensure everything works:
```bash
go test ./...
```

4. Start a local SpiceDB for testing:
```bash
docker run --rm -p 50051:50051 authzed/spicedb serve \
  --grpc-preshared-key "devkey" \
  --datastore-engine memory
```

## Development Workflow

### Code Style

- Follow Go conventions and idioms
- Use `gofmt` for formatting: `gofmt -s -w .`
- Use `go vet` to check for common errors: `go vet ./...`
- Document public functions with godoc comments

### Testing

- Write tests for new features
- Ensure all tests pass: `go test ./...`
- Aim for good code coverage

Running specific tests:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./spicedb

# Run a specific test
go test -run TestCheckPermission ./spicedb
```

### Running Examples

```bash
# Set up environment
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
export SPICEDB_INSECURE=true

# Run basic examples
go run ./examples/basic/main.go

# Run HTTP service
go run ./examples/http-service/main.go
```

## Making Changes

### Creating a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### Making Your Changes

1. Make your changes to the code
2. Add or update tests
3. Update documentation as needed
4. Run tests and linters

### Common Development Tasks

#### Adding a New Client Method

1. Add the method to `spicedb/client.go`
2. Create a corresponding builder in `spicedb/builders.go` if applicable
3. Add tests in `spicedb/client_test.go`
4. Update documentation in `spicedb/README.md`

#### Adding Error Types

1. Add to `spicedb/errors.go`
2. Create helper functions to check error type
3. Add tests and documentation

#### Adding Examples

1. Create example file in `examples/` directory
2. Include in `examples/README.md`
4. Test the example works

## Commit Messages

Use clear, descriptive commit messages:

```
feat: add support for relationship lookups

Implement LookupResources and LookupSubjects operations
with full builder pattern support.

Fixes #123
```

Guidelines:
- Use imperative mood ("add" not "added")
- First line should be 50 characters or less
- Reference issues when applicable
- Explain what and why, not how

## Pull Requests

### Before Submitting

1. Ensure all tests pass: `go test ./...`
2. Run linters: `go vet ./...` and `gofmt -l .`
3. Update documentation
4. Squash commits if needed

### PR Template

```
## Description

Brief description of what this PR does.

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Related Issues

Closes #(issue number)

## Testing

- [ ] Added tests
- [ ] All tests pass
- [ ] Tested manually

## Documentation

- [ ] Updated README
- [ ] Updated godoc comments
- [ ] Updated examples

## Breaking Changes

If applicable, describe any breaking changes.
```

## Code Review Process

- At least one maintainer review required
- All tests must pass
- No unresolved conversations
- Changes must not decrease code quality or test coverage

## Documentation Guidelines

### Godoc Comments

```go
// CheckPermission checks if a subject has permission on a resource.
// It returns the permission status and any errors that occurred.
//
// The request is cached for 5 seconds by default. Sets with context
// caveats require the caveat condition to be satisfied.
//
// Example:
//   resp, err := client.CheckPermission(ctx, CheckRequest{
//     Subject:    "user:alice",
//     Resource:   "doc:1",
//     Permission: "view",
//   })
//
// See ClientOptions for configuration options.
func (c *Client) CheckPermission(ctx context.Context, req CheckRequest) (CheckResponse, error) {
```

### README Updates

- Keep examples up to date
- Explain why, not just how
- Include links to related resources
- Use consistent formatting

## Reporting Issues

### Bug Reports

Include:
- Go version
- Library version
- Minimal code to reproduce
- Expected vs actual behavior
- Error messages or logs

### Feature Requests

Include:
- Use case description
- Proposed API design
- Alternative approaches considered

## Questions?

- Check existing issues and documentation
- Start a discussion in GitHub Discussions
- Join the [AuthZed community](https://authzed.com/community)

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.

## Recognition

We appreciate all contributions! Contributors will be recognized in appropriate places.

---

Thank you for contributing to SpiceDB Go Library! 🎉

