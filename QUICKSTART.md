# Quick Start Guide

Get the SpiceDB Go Library up and running in 5 minutes!

## Prerequisites

- Go 1.20 or later
- Docker (optional, for local SpiceDB)
- Access to a SpiceDB instance

## Installation

```bash
go get github.com/Yasiruofficial/events-authz
```

## 1. Start SpiceDB (Development)

```bash
docker run --rm -p 50051:50051 authzed/spicedb serve \
  --grpc-preshared-key "devkey" \
  --datastore-engine memory
```

## 2. Create Your First Application

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
	// Create client
	client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Check permission
	resp, err := client.CheckPermission(context.Background(), types.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:budget-2026",
		Permission: "view",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Permission allowed: %v", resp.Allowed)
}
```

### Using Builder Pattern (Recommended)

```go
allowed, err := client.CheckPermissionBuilder().
	Subject("user:bob").
	Resource("document:report").
	Permission("edit").
	IsAllowed(context.Background())
if err != nil {
	log.Fatal(err)
}
log.Printf("Bob can edit: %v", allowed)
```

## 3. Common Patterns

### With Consistency Control

```go
resp, err := client.CheckPermission(ctx, types.CheckRequest{
	Subject:     "user:alice",
	Resource:    "account:ACC-123",
	Permission:  "debit",
	Consistency: "fully_consistent", // Always get latest data
})
```

### With Caveat Context

```go
resp, err := client.CheckPermission(ctx, types.CheckRequest{
	Subject:    "user:charlie",
	Resource:   "file:confidential.pdf",
	Permission: "download",
	Context: map[string]interface{}{
		"ip_address": "203.0.113.42",
		"time":       "business_hours",
	},
})
```

### Batch Checks

```go
resources := []string{"doc:1", "doc:2", "doc:3"}
for _, resource := range resources {
	resp, err := client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:alice",
		Resource:   resource,
		Permission: "view",
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("%s: %v", resource, resp.Allowed)
	}
}
```

## 4. Error Handling

```go
resp, err := client.CheckPermission(ctx, request)
if err != nil {
	if spicedb.IsValidationError(err) {
		log.Printf("Bad request: %v", err)
		// Handle validation error
	} else if spicedb.IsOperationError(err) {
		log.Printf("SpiceDB error: %v", err)
		// Handle operation error
	} else {
		log.Printf("Unknown error: %v", err)
	}
}
```

## 5. Configuration

### Development (Default)

```go
client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
```

### Production

```go
client, err := spicedb.NewClient(spicedb.ClientOptions{
	Address:            "spicedb.prod.example.com:50051",
	PreSharedKey:       os.Getenv("SPICEDB_TOKEN"),
	TLSEnabled:         true,
	InsecureSkipVerify: false,
	RequestTimeout:     5 * time.Second,
	DefaultConsistency: "fully_consistent",
})
```

### With Custom Cache

```go
client, err := spicedb.NewClient(spicedb.ClientOptions{
	Address:      "localhost:50051",
	PreSharedKey: "devkey",
	Cache:        myCustomCache, // Implement cache.Interface
})
```

## 6. Testing

### Unit Tests

```bash
go test ./spicedb
```

### Integration Tests

```bash
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
go test ./spicedb -run Integration
```

## 7. Examples

Run comprehensive examples:

```bash
# Basic usage examples
go run ./examples/basic/main.go

# HTTP service wrapper
go run ./examples/http-service/main.go
```

## 8. Environment Variables

For examples and HTTP service:

```bash
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
export SPICEDB_INSECURE=true
export SPICEDB_CONSISTENCY=minimize_latency
export SPICEDB_TIMEOUT=3s
export HTTP_ADDR=:8080
```

## 9. Next Steps

- Read [spicedb/README.md](./spicedb/README.md) for full API documentation
- Check [examples/](./examples/) for more examples
- See [CONTRIBUTING.md](./CONTRIBUTING.md) for development guidelines
- Review [ROADMAP.md](./ROADMAP.md) for planned features

## 10. Common Issues

### Connection Error

```
Error: connection refused
```

**Solution**: Make sure SpiceDB is running on the specified address:

```bash
docker run --rm -p 50051:50051 authzed/spicedb serve \
  --grpc-preshared-key "devkey" \
  --datastore-engine memory
```

### Validation Error

```
Error: invalid request on field 'subject': must be in type:id format
```

**Solution**: Ensure resource and subject follow the format `type:id`:

```go
// ✓ Correct
subject: "user:alice"
resource: "document:budget-2026"

// ✗ Wrong
subject: "alice"
resource: "budget-2026"
```

### Permission Denied

*This is often expected!* It usually means no permission was granted, which is not an error:

```go
resp, err := client.CheckPermission(ctx, request)
if err != nil {
	// Error with the request itself
	log.Fatal(err)
}
// resp.Allowed will be false if permission is not granted
log.Printf("Access: %v", resp.Allowed)
```

## 11. Debugging

### Enable Verbose Logging

Most Go libraries can be debugged with environment variables:

```bash
# For gRPC debugging
export GRPC_GO_LOG_VERBOSITY_LEVEL=99
export GRPC_GO_LOG_SEVERITY_LEVEL=info
```

### Test Connection

Simple test to verify server connectivity:

```go
client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
if err != nil {
	log.Printf("Failed to connect: %v", err)
	// Debug connection issues
}
defer client.Close()
log.Println("Connected to SpiceDB!")
```

## 12. Performance Tips

1. **Enable Caching**: Default in-memory cache can significantly improve performance
2. **Use minimize_latency**: Faster response times for non-critical checks
3. **Connection Reuse**: Keep client open, don't create new ones per request
4. **Batch Requests**: Group multiple checks into a single operation when possible
5. **Context Timeout**: Set reasonable timeouts to prevent hanging

```go
// Good: Reuse client
client, _ := spicedb.NewClientWithDefaults(addr, token)
defer client.Close()

// Handle many requests with the same client
for req := range requests {
	resp, _ := client.CheckPermission(ctx, req)
	// process response
}
```

## Resources

- [Full Documentation](./spicedb/README.md)
- [SpiceDB Docs](https://authzed.com/docs)
- [AuthZed Community](https://authzed.com/community)
- [Examples](./examples/)

---

**Happy coding!** 🚀

