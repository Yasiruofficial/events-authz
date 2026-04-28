# Migration Guide: HTTP Service to Library

This guide helps you migrate from the old HTTP service to the new SpiceDB library.

## Overview

**Old Structure**: Monolithic HTTP service with embedded authorization logic

**New Structure**: Reusable library + optional HTTP example

## Two Migration Paths

### Path 1: Use the Library Directly (Recommended)

Import the library into your application:

```go
import "github.com/Yasiruofficial/events-authz/spicedb"
```

**Benefits:**
- No HTTP overhead
- Direct control over authorization flow
- Better integration with existing code
- Smaller codebase for embedding

**Example:**

```go
// Old: HTTP service as standalone
// Now: Import and use directly in your app
client, _ := spicedb.NewClientWithDefaults(addr, token)
resp, _ := client.CheckPermission(ctx, request)
```

### Path 2: Use HTTP Example (Migration Path)

Re-deploy the HTTP service example:

```bash
go run ./examples/http-service/main.go
```

**Benefits:**
- Familiar HTTP interface
- No code changes needed
- Gradual migration to library
- Good for services behind API gateway

**Same API as before:**

```bash
POST /check
{
  "subject": "user:alice",
  "resource": "document:budget",
  "permission": "view"
}
```

## Step-by-Step Migration

### Step 1: Investigate Current Usage

Check how the old service is being used:

```bash
# If using HTTP API
grep -r "http://.*:8080/check" . 

# If using as library
grep -r "spicedb.NewClient" .
```

### Step 2: Choose Migration Path

**Choose Direct Library if:**
- You control the codebase
- Can update imports
- Want better performance
- Writing new code

**Choose HTTP Example if:**
- Using via HTTP API
- External consumers
- Quick migration needed
- Want to keep interface stable

### Step 3: Update Configuration

#### Old Configuration
```bash
export HTTP_ADDR=:8080
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
```

#### New Configuration (Library)
```bash
# In code:
client, err := spicedb.NewClient(spicedb.ClientOptions{
    Address:      "localhost:50051",
    PreSharedKey: "devkey",
})

// Or environment variables (examples only):
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
```

#### New Configuration (HTTP Example)
```bash
# Same as before!
export HTTP_ADDR=:8080
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
go run ./examples/http-service/main.go
```

### Step 4: Update Code

#### Option A: Direct Library Usage

**Before (HTTP):**
```go
import (
    "events-authz/internal/http"
    "events-authz/internal/service"
    "events-authz/internal/spicedb"
)

func main() {
    spice, _ := spicedb.NewClient(...)
    svc := service.NewAuthzService(spice, cache)
    handler := http.NewHandler(svc)
    // Start HTTP server
}
```

**After (Library):**
```go
import "github.com/Yasiruofficial/events-authz/spicedb"
import "github.com/Yasiruofficial/events-authz/spicedb/types"

func main() {
    client, _ := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
    defer client.Close()
    
    resp, _ := client.CheckPermission(ctx, types.CheckRequest{
        Subject:    "user:alice",
        Resource:   "document:1",
        Permission: "view",
    })
}
```

#### Option B: HTTP Service (Minimal Changes)

Run the example service:

```bash
go run ./examples/http-service/main.go
```

The API remains the same, so no client code changes needed.

### Step 5: Testing

#### Old Tests
```go
// Old test setup
client, _ := spicedb.NewClient(addr, token, insecure, timeout, consistency)
```

#### New Tests
```go
// New test setup
client, _ := spicedb.NewClient(spicedb.ClientOptions{
    Address:      addr,
    PreSharedKey: token,
    TLSEnabled:   !insecure,
    RequestTimeout: timeout,
    DefaultConsistency: consistency,
})
```

**Example Test:**
```go
func TestPermissionCheck(t *testing.T) {
    client, _ := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
    defer client.Close()
    
    resp, err := client.CheckPermission(context.Background(), types.CheckRequest{
        Subject:    "user:test",
        Resource:   "doc:test",
        Permission: "view",
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
}
```

### Step 6: Verification

Verify migration worked:

#### Library Usage
```bash
go test ./...
go build ./...
```

#### HTTP Service Usage
```bash
go run ./examples/http-service/main.go &
curl -X POST http://localhost:8080/check \
  -H "Content-Type: application/json" \
  -d '{"subject":"user:test","resource":"doc:test","permission":"view"}'
```

## API Comparison

### Permission Check

#### Old
```go
response, err := spice.CheckPermission(ctx, model.CheckRequest{
    Subject:     "user:alice",
    Resource:    "document:1",
    Permission: "view",
    Context:    map[string]any{...},
    Consistency: "fully_consistent",
    ZedToken:   "token",
})
```

#### New - Direct
```go
response, err := client.CheckPermission(ctx, types.CheckRequest{
    Subject:     "user:alice",
    Resource:    "document:1",
    Permission: "view",
    Context:    map[string]interface{}{...},
    Consistency: "fully_consistent",
    ZedToken:   "token",
})
```

#### New - Builder
```go
response, err := client.CheckPermissionBuilder().
    Subject("user:alice").
    Resource("document:1").
    Permission("view").
    WithContext(map[string]interface{}{...}).
    WithConsistency("fully_consistent").
    WithZedToken("token").
    Check(ctx)
```

### Error Handling

#### Old
```go
if errors.Is(err, spicedb.ErrInvalidCheckRequest) {
    // validation error
}
```

#### New
```go
if spicedb.IsValidationError(err) {
    // validation error
} else if spicedb.IsOperationError(err) {
    // operation error
}
```

## Breaking Changes

1. **Module Path**: `events-authz` → `github.com/Yasiruofficial/events-authz`
2. **Package Structure**: `internal/` packages → public `spicedb/` package
3. **Type Names**: `model.CheckRequest` → `types.CheckRequest`
4. **Error Handling**: Custom checks → `IsValidationError()`, `IsOperationError()`
5. **Options Structure**: Direct parameters → `ClientOptions` struct

## Gradual Migration Strategy

If you have multiple services:

### Phase 1: Deploy HTTP Example
- Run HTTP service example alongside old service
- Update consumers to use new endpoint
- Verify functionality

### Phase 2: Migrate to Direct Library
- Update internal services to use library directly
- Remove HTTP wrapper
- Performance improvement

### Phase 3: Cleanup
- Deprecate old code
- Archive old implementation
- Full adoption of library

## Performance Implications

### Library (Direct)
- ✅ Faster: No HTTP overhead
- ✅ Lower memory: No web framework
- ✅ Better caching: In-process cache

### HTTP Service
- ✔️ Same latency: Network vs process
- ⚠️ More overhead: HTTP serialization
- ✔️ Better isolation: Separate service

## Rollback Plan

If issues occur:

1. **HTTP Example**: Just restart old service
2. **Library**: Revert code changes, keep old branch available

```bash
# Keep old code
git branch v1-old-service
git checkout v1-old-service
# Restart old service
```

## Common Migration Issues

### Issue: Import Paths

**Error**: `cannot find package "events-authz/internal/spicedb"`

**Solution**: Update imports:
```go
// Old
import "events-authz/internal/spicedb"

// New
import "github.com/Yasiruofficial/events-authz/spicedb"
```

### Issue: Type Mismatches

**Error**: `model.CheckRequest` not assignable to `types.CheckRequest`

**Solution**: Update type references:
```go
// Old
var req model.CheckRequest

// New
var req types.CheckRequest
```

### Issue: Context Type

**Error**: `map[string]any` vs `map[string]interface{}`

**Solution**: Both work, but new code uses `interface{}`:
```go
Context: map[string]interface{}{...}
```

## Support

- **Questions**: Open a GitHub issue
- **Problems**: Check [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)
- **Examples**: See [examples/](./examples/)
- **Full Docs**: See [spicedb/README.md](./spicedb/README.md)

## Checklist

Ensure these steps are complete:

- [ ] Choose migration path
- [ ] Update imports
- [ ] Update configuration
- [ ] Update code/tests
- [ ] Run tests
- [ ] Deploy and verify
- [ ] Monitor for issues
- [ ] Archive old code (optional)

---

**For questions or issues, see [CONTRIBUTING.md](./CONTRIBUTING.md)**

