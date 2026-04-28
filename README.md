# SpiceDB Go Library

[![Go Reference](https://pkg.go.dev/badge/github.com/Yasiruofficial/events-authz.svg)](https://pkg.go.dev/github.com/Yasiruofficial/events-authz)
[![Go Report Card](https://goreportcard.com/badge/github.com/Yasiruofficial/events-authz)](https://goreportcard.com/report/github.com/Yasiruofficial/events-authz)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A production-ready Go library for interacting with [SpiceDB](https://authzed.com/spicedb) - an open-source authorization system. This library provides a clean, idiomatic Go API for permission checks, relationship management, and resource lookups.

## ✨ Features

- **Permission Checks**: Check if a subject has permission on a resource
- **Relationship Management**: Create, read, update, and delete relationships
- **Resource Lookups**: Find resources that a subject has access to
- **Subject Lookups**: Find subjects that have access to a resource
- **Flexible Consistency**: Support for multiple consistency levels
- **Built-in Caching**: Optional configurable caching layer
- **Type-Safe**: Comprehensive error handling with typed errors
- **Builder Pattern**: Fluent API for building complex requests
- **Full Context Support**: Go context support for cancellation and timeouts
- **Production Ready**: Used in production environments
- **Well Documented**: Comprehensive documentation and examples

## 📦 Installation

```bash
go get github.com/Yasiruofficial/events-authz
```

## 🚀 Quick Start

### Basic Permission Check

```go
package main

import (
	"context"
	"log"

	"github.com/Yasiruofficial/events-authz/spicedb"
	"github.com/Yasiruofficial/events-authz/spicedb/types"
)

func main() {
	// Create a client
	client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Check permission
	response, err := client.CheckPermission(context.Background(), types.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:budget-2026",
		Permission: "view",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Permission allowed: %v", response.Allowed)
}
```

### Using the Builder Pattern

```go
allowed, err := client.CheckPermissionBuilder().
	Subject("user:alice").
	Resource("document:budget-2026").
	Permission("view").
	IsAllowed(context.Background())
if err != nil {
	log.Fatal(err)
}
log.Printf("Allowed: %v", allowed)
```

## 📚 Documentation

- **[Library Documentation](./spicedb/README.md)** - Comprehensive guide and API reference
- **[Examples](./examples/)** - Sample code demonstrating various features
- **[HTTP Service Example](./examples/http-service/)** - Complete HTTP wrapper service

## 🏗️ Project Structure

```
spicedb-go/
├── spicedb/              # Main library package
│   ├── README.md         # Full documentation
│   ├── client.go         # Main client implementation
│   ├── builders.go       # Builder patterns
│   ├── errors.go         # Error types
│   ├── types/            # Type definitions
│   ├── cache/            # Caching abstraction
│   └── client_test.go    # Tests
├── examples/             # Example applications
│   ├── http-service/     # HTTP API wrapper
│   ├── basic/            # Basic usage examples
│   └── README.md         # Examples guide
├── go.mod                # Module definition
├── go.sum                # Dependencies checksums
└── README.md             # This file
```

## 🔧 Configuration

### With Defaults (Development)

```go
client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
```

### With Custom Options

```go
client, err := spicedb.NewClient(spicedb.ClientOptions{
	Address:             "spicedb.example.com:50051",
	PreSharedKey:        "your-secret-key",
	TLSEnabled:          true,
	InsecureSkipVerify:  false,
	RequestTimeout:      5 * time.Second,
	DefaultConsistency:  "fully_consistent",
	DisableCache:        false,
})
```

## 🧪 Environment Variables (For Examples)

- `SPICEDB_ADDR`: SpiceDB server address (default: `localhost:50051`)
- `SPICEDB_TOKEN`: Pre-shared key for authentication
- `SPICEDB_INSECURE`: Use insecure connection (default: `true`)
- `SPICEDB_CONSISTENCY`: Default consistency level (default: `minimize_latency`)
- `SPICEDB_TIMEOUT`: Request timeout (default: `3s`)
- `HTTP_ADDR`: HTTP service address (default: `:8080`) - for HTTP service example

## 🚀 Getting Started

### 1. Start a local SpiceDB instance

```bash
docker run --rm -p 50051:50051 authzed/spicedb serve \
  --grpc-preshared-key "devkey" \
  --datastore-engine memory
```

### 2. Run an example

```bash
cd examples/basic
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
export SPICEDB_INSECURE=true
go run main.go
```

### 3. Or use the HTTP service example

```bash
github.com/Yasiruofficial/events-authz/examples/http-service
SPICEDB_ADDR=localhost:50051 \
SPICEDB_TOKEN=devkey \
SPICEDB_INSECURE=true \
go run main.go
```

Then test it:

```bash
curl -X POST http://localhost:8080/check \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "user:alice",
    "resource": "document:budget-2026",
    "permission": "view"
  }'
```

## 📖 Usage Examples

### Permission Check with Context (Caveats)

```go
response, err := client.CheckPermission(ctx, types.CheckRequest{
	Subject:    "user:bob",
	Resource:   "file:confidential.pdf",
	Permission: "download",
	Context: map[string]interface{}{
		"ip_address": "203.0.113.42",
		"time":       "2024-01-15T14:30:00Z",
	},
})
```

### Batch Permission Checks

```go
requests := []types.CheckRequest{
	{Subject: "user:alice", Resource: "doc:1", Permission: "view"},
	{Subject: "user:bob", Resource: "doc:2", Permission: "edit"},
}

for _, req := range requests {
	resp, err := client.CheckPermission(ctx, req)
	// handle response
}
```

### Custom Consistency Level

```go
allowed, err := client.CheckPermissionBuilder().
	Subject("user:alice").
	Resource("account:ACC-123").
	Permission: "view").
	WithConsistency("fully_consistent").
	IsAllowed(ctx)
```

### Error Handling

```go
resp, err := client.CheckPermission(ctx, request)
if err != nil {
	if spicedb.IsValidationError(err) {
		// Request validation failed
		log.Printf("Validation error: %v", err)
	} else if spicedb.IsOperationError(err) {
		// SpiceDB operation failed
		log.Printf("Operation error: %v", err)
	}
}
```

## 🔐 Requirements

- Go 1.20 or later
- Access to a running SpiceDB instance
- For development: Docker (to run SpiceDB locally)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Resources

- [SpiceDB Documentation](https://authzed.com/docs)
- [AuthZed Community](https://authzed.com/community)
- [SpiceDB GitHub](https://github.com/authzed/spicedb)
- [AuthZed Go Client](https://github.com/authzed/authzed-go)

## 📞 Support

- **Issues**: Report bugs on [GitHub Issues](https://github.com/Yasiruofficial/events-authz/issues)
- **Discussions**: Join community discussions
- **Documentation**: See [./spicedb/README.md](./spicedb/README.md)

---

**Built with ❤️ by SpiceDB Community**
