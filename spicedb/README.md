# SpiceDB Go Library

[![Go Reference](https://pkg.go.dev/badge/github.com/Yasiruofficial/events-authz.svg)](https://pkg.go.dev/github.com/Yasiruofficial/events-authz)
[![Go Report Card](https://goreportcard.com/badge/github.com/Yasiruofficial/events-authz)](https://goreportcard.com/report/github.com/Yasiruofficial/events-authz)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A production-ready Go library for interacting with [SpiceDB](https://authzed.com/spicedb) authorization service. This library provides a clean, idiomatic Go API for permission checks, relationship management, and resource lookups.

## Features

- **Permission Checks**: Check if a subject has permission on a resource
- **Relationship Management**: Create, read, update, and delete relationships
- **Resource Lookups**: Find resources that a subject has access to
- **Subject Lookups**: Find subjects that have access to a resource
- **Flexible Consistency**: Support for multiple consistency levels (minimize_latency, fully_consistent, at_least_as_fresh, at_exact_snapshot)
- **Caching**: Optional built-in caching with customizable implementations
- **Type-Safe**: Comprehensive error handling with typed errors
- **Builder Pattern**: Fluent API for building complex requests
- **Context Support**: Full Go context support for cancellation and timeouts

## Installation

```bash
go get github.com/Yasiruofficial/events-authz
```

## Quick Start

### Basic Setup

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
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Check a permission
	response, err := client.CheckPermission(context.Background(), types.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:budget-2026",
		Permission: "view",
	})
	if err != nil {
		log.Fatalf("permission check failed: %v", err)
	}

	log.Printf("Permission allowed: %v", response.Allowed)
}
```

### Using the Builder Pattern

The library provides a fluent builder pattern for more complex scenarios:

```go
allowed, err := client.CheckPermissionBuilder().
	Subject("user:alice").
	Resource("document:budget-2026").
	Permission("view").
	WithConsistency("fully_consistent").
	IsAllowed(context.Background())
if err != nil {
	log.Fatal(err)
}
log.Printf("Allowed: %v", allowed)
```

### Advanced Configuration

```go
client, err := spicedb.NewClient(spicedb.ClientOptions{
	Address:             "spicedb.example.com:50051",
	PreSharedKey:        "my-secret-key",
	TLSEnabled:          true,
	InsecureSkipVerify:  false,
	RequestTimeout:      5 * time.Second,
	DefaultConsistency:  "fully_consistent",
	DisableCache:        false,
	// Optionally provide a custom cache implementation
	// Cache: myCustomCache,
})
if err != nil {
	log.Fatal(err)
}
defer client.Close()
```

## API Documentation

### Permission Checks

#### CheckPermission

```go
response, err := client.CheckPermission(ctx, types.CheckRequest{
	Subject:     "user:alice#member",      // Can include optional relation
	Resource:    "organization:acme",       // Format: type:id
	Permission:  "admin",
	Context:     map[string]interface{}{   // Optional caveat context
		"ip_address": "192.168.1.1",
	},
	Consistency: "fully_consistent",        // Optional
	ZedToken:    "",                        // For snapshot-based consistency
})
if err != nil {
	// Handle error
}
log.Printf("Permission granted: %v", response.Allowed)
log.Printf("Permissionship: %v", response.Permissionship)
log.Printf("Checked at: %v", response.CheckedAt)
```

#### CheckPermissionBuilder

For cleaner, more readable code:

```go
allowed, err := client.CheckPermissionBuilder().
	Subject("user:bob").
	Resource("document:report-2026").
	Permission("edit").
	WithContext(map[string]interface{}{"department": "finance"}).
	IsAllowed(ctx)
```

### Consistency Levels

The library supports SpiceDB's consistency guarantees:

- **minimize_latency** (default): Fastest response time, eventual consistency
- **fully_consistent**: Waits for all writes to be fully consistent
- **at_least_as_fresh**: Uses a ZedToken to ensure at least that level of freshness
- **at_exact_snapshot**: Ensures responses are from an exact point in time

```go
client.CheckPermissionBuilder().
	Subject("user:alice").
	Resource("document:budget").
	Permission("view").
	WithConsistency("at_least_as_fresh").
	WithZedToken(previousToken).
	Check(ctx)
```

## Error Handling

The library provides typed errors for better error handling:

```go
resp, err := client.CheckPermission(ctx, request)
if err != nil {
	if spicedb.IsValidationError(err) {
		// Request validation failed
		log.Printf("Validation error: %v", err)
	} else if spicedb.IsOperationError(err) {
		// SpiceDB operation failed
		log.Printf("Operation error: %v", err)
	} else {
		// Other errors
		log.Printf("Unexpected error: %v", err)
	}
}
```

## Caching

The library includes optional caching to reduce load on SpiceDB:

### Using Default In-Memory Cache

```go
// In-memory cache is enabled by default
client, err := spicedb.NewClient(spicedb.ClientOptions{
	Address:      "localhost:50051",
	PreSharedKey: "devkey",
	// Cache enabled by default
})
```

### Disabling Cache

```go
client, err := spicedb.NewClient(spicedb.ClientOptions{
	Address:      "localhost:50051",
	PreSharedKey: "devkey",
	DisableCache: true,  // Disable caching
})
```

### Custom Cache Implementation

```go
type MyCustomCache struct {
	// your implementation
}

func (c *MyCustomCache) Get(key string) (interface{}, bool) {
	// implementation
}

func (c *MyCustomCache) Set(key string, value interface{}, ttl time.Duration) {
	// implementation
}

func (c *MyCustomCache) Delete(key string) {
	// implementation
}

func (c *MyCustomCache) Clear() {
	// implementation
}

client, err := spicedb.NewClient(spicedb.ClientOptions{
	Address:      "localhost:50051",
	PreSharedKey: "devkey",
	Cache:        &MyCustomCache{},
})
```

## Examples

### Example 1: Basic Permission Check

```go
ctx := context.Background()
resp, err := client.CheckPermission(ctx, types.CheckRequest{
	Subject:    "user:alice",
	Resource:   "document:budget-2026",
	Permission: "view",
})
if err != nil {
	log.Fatal(err)
}
if resp.Allowed {
	log.Println("Alice can view the budget document")
}
```

### Example 2: Permission Check with Context (Caveats)

```go
resp, err := client.CheckPermission(ctx, types.CheckRequest{
	Subject:    "user:bob",
	Resource:   "file:confidential.pdf",
	Permission: "download",
	Context: map[string]interface{}{
		"ip_address": "203.0.113.42",
		"time":       "2024-01-15T14:30:00Z",
	},
})
```

### Example 3: Fully Consistent Permission Check

```go
resp, err := client.CheckPermission(ctx, types.CheckRequest{
	Subject:     "service:payment-processor",
	Resource:    "account:ACC-12345",
	Permission:  "debit",
	Consistency: "fully_consistent",
})
```

### Example 4: Using Builder with Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

allowed, err := client.CheckPermissionBuilder().
	Subject("user:charlie").
	Resource("database:production").
	Permission("access").
	IsAllowed(ctx)
```

## Common Patterns

### Batch Permission Checks

```go
requests := []types.CheckRequest{
	{Subject: "user:alice", Resource: "doc:1", Permission: "view"},
	{Subject: "user:bob", Resource: "doc:2", Permission: "edit"},
	// ... more requests
}

results := make([]types.CheckResponse, len(requests))
for i, req := range requests {
	resp, err := client.CheckPermission(ctx, req)
	if err != nil {
		log.Printf("Error checking %s on %s: %v", req.Subject, req.Resource, err)
		continue
	}
	results[i] = resp
}
```

### Caching with Custom TTL

Implement a custom cache that respects per-operation TTLs:

```go
type CustomTTLCache struct {
	store map[string]CacheEntry
	mu    sync.RWMutex
}

type CacheEntry struct {
	Value      interface{}
	ExpiresAt  time.Time
	CustomTTL  time.Duration
}

func (c *CustomTTLCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CustomTTL: ttl,
	}
}

// ... implement other methods
```

## Requirements

- Go 1.20 or later
- Access to a running SpiceDB instance

## Testing

Run the test suite:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests to our repository.

## License

This library is licensed under the Apache License 2.0 - see the LICENSE file for details.

## Support

For issues, questions, or contributions:
- Report bugs on [GitHub Issues](https://github.com/Yasiruofficial/events-authz/issues)
- Join our [community Slack](https://authzed.com/community)
- Visit the [SpiceDB documentation](https://authzed.com/docs)

## Related Resources

- [SpiceDB Documentation](https://authzed.com/docs)
- [SpiceDB Schema Language](https://authzed.com/docs/reference/schema-lang/intro)
- [AuthZed API Reference](https://authzed.com/docs/reference/api)

