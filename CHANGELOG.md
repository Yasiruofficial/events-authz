# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial public release of SpiceDB Go Library
- Core `Client` with CheckPermission support
- Builder pattern for fluent API
- In-memory caching implementation
- Type-safe error handling
- Comprehensive documentation and examples
- HTTP service example application
- Support for all SpiceDB consistency levels
- Context and caveat support

### Changed
- Refactored from HTTP service to reusable library

### Removed
- Direct HTTP binding from core library (moved to examples)
- Monolithic server structure

## [0.1.0-alpha] - 2024-01-15

### Added
- Initial release
- CheckPermission operation
- Builder pattern for requests
- In-memory cache with TTL
- Error types and validation
- HTTP wrapper example
- Full documentation

---

## Release Guidelines

### Versioning

This project follows [Semantic Versioning](https://semver.org/):
- MAJOR version for incompatible API changes
- MINOR version for backwards-compatible functionality additions
- PATCH version for backwards-compatible bug fixes

### Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create git tag (e.g., `v0.2.0`)
4. Push changes and tags
5. Create GitHub release with changes

### Support Windows

- Latest version only (N)
- Previous version (N-1) for critical bugs only
- Older versions archived/deprecated with announcement

## Deprecated Features

Currently, all features are new and none are deprecated.

---

For a list of planned features, see [ROADMAP.md](ROADMAP.md).

