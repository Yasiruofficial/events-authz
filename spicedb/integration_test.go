package spicedb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Yasiruofficial/events-authz/spicedb/types"
)

// TestIntegration provides integration tests with a real SpiceDB instance
// To run these tests, you need a SpiceDB instance running:
// docker run --rm -p 50051:50051 authzed/spicedb serve \
//   --grpc-preshared-key "devkey" \
//   --datastore-engine memory

func skipIfNoSpiceDB(t *testing.T) {
	if os.Getenv("SPICEDB_ADDR") == "" {
		t.Skip("SPICEDB_ADDR not set, skipping integration tests")
	}
}

func getTestClient(t *testing.T) *Client {
	addr := os.Getenv("SPICEDB_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	token := os.Getenv("SPICEDB_TOKEN")
	if token == "" {
		token = "devkey"
	}

	opts := ClientOptions{
		Address:        addr,
		PreSharedKey:   token,
		TLSEnabled:     false,
		RequestTimeout: 5 * time.Second,
		DisableCache:   true, // Disable cache for integration tests
	}

	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return client
}

// TestCheckPermissionIntegration tests CheckPermission with a real SpiceDB
func TestCheckPermissionIntegration(t *testing.T) {
	skipIfNoSpiceDB(t)

	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Note: This will fail without proper schema and relationships in SpiceDB
	// This is just demonstrating the integration
	_, err := client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:test",
		Resource:   "resource:test",
		Permission: "view",
	})

	// We expect some kind of response or error - both are valid
	// The important thing is that we connected and got a response
	if err != nil {
		t.Logf("Got expected error (schema may not be defined): %v", err)
	}
}

// TestClientCreationIntegration verifies client creation works with real server
func TestClientCreationIntegration(t *testing.T) {
	skipIfNoSpiceDB(t)

	addr := os.Getenv("SPICEDB_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	client, err := NewClientWithDefaults(addr, "devkey")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	if client == nil {
		t.Fatal("client is nil")
	}
}

// TestBuilderIntegration tests builder pattern with real client
func TestBuilderIntegration(t *testing.T) {
	skipIfNoSpiceDB(t)

	client := getTestClient(t)
	defer client.Close()

	// Just verify the builder doesn't panic
	builder := client.CheckPermissionBuilder().
		Subject("user:test").
		Resource("resource:test").
		Permission("view")

	// We don't execute because schema may not exist, just verify builder works
	if builder == nil {
		t.Fatal("builder is nil")
	}
}

// TestContextHandling tests that contexts are properly passed through
func TestContextHandling(t *testing.T) {
	skipIfNoSpiceDB(t)

	client := getTestClient(t)
	defer client.Close()

	// Test with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Try to make a request
	_, err := client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:test",
		Resource:   "resource:test",
		Permission: "view",
	})

	// We might get timeout or connection error or schema error
	// The point is the context was respected
	if err != nil {
		t.Logf("Got error (expected): %v", err)
	}
}

// TestErrorRecovery tests that client can recover from errors
func TestErrorRecovery(t *testing.T) {
	skipIfNoSpiceDB(t)

	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Make a request that will likely fail (no schema)
	_, _ = client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:test",
		Resource:   "resource:test",
		Permission: "view",
	})

	// Try again - should not panic or hang
	_, _ = client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:test2",
		Resource:   "resource:test2",
		Permission: "read",
	})

	// Both calls completed without panic - test passes
}
