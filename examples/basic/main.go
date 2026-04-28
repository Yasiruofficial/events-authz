package main

import (
	"context"
	"log"
	"time"

	"github.com/Yasiruofficial/events-authz/spicedb"
	"github.com/Yasiruofficial/events-authz/spicedb/types"
)

// BasicExample demonstrates basic permission checks
func BasicExample() {
	client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Check if alice can view budget-2026
	response, err := client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:budget-2026",
		Permission: "view",
	})
	if err != nil {
		log.Fatalf("permission check failed: %v", err)
	}

	log.Printf("Alice can view document: %v", response.Allowed)
	log.Printf("Permissionship: %s", response.Permissionship)
	log.Printf("Checked at: %s", response.CheckedAt)
}

// BuilderExample demonstrates using the fluent builder pattern
func BuilderExample() {
	client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Using builder pattern for cleaner code
	allowed, err := client.CheckPermissionBuilder().
		Subject("user:bob").
		Resource("document:report-2026").
		Permission("edit").
		IsAllowed(ctx)
	if err != nil {
		log.Fatalf("permission check failed: %v", err)
	}

	log.Printf("Bob can edit report: %v", allowed)
}

// ContextExample demonstrates permission checks with context (caveats)
func ContextExample() {
	client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Check with additional context for caveats
	response, err := client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:charlie",
		Resource:   "file:confidential.pdf",
		Permission: "download",
		Context: map[string]interface{}{
			"ip_address":  "203.0.113.42",
			"time_of_day": "business_hours",
			"department":  "finance",
		},
	})
	if err != nil {
		log.Fatalf("permission check failed: %v", err)
	}

	log.Printf("Charlie can download file: %v", response.Allowed)
}

// ConsistencyExample demonstrates different consistency levels
func ConsistencyExample() {
	client, err := spicedb.NewClient(spicedb.ClientOptions{
		Address:            "localhost:50051",
		PreSharedKey:       "devkey",
		TLSEnabled:         false,
		DefaultConsistency: "minimize_latency",
		RequestTimeout:     5 * time.Second,
	})
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Fast check with eventual consistency (default)
	fastResp, _ := client.CheckPermission(ctx, types.CheckRequest{
		Subject:     "user:alice",
		Resource:    "account:ACC-123",
		Permission:  "view",
		Consistency: "minimize_latency",
	})
	log.Printf("Fast check result: %v", fastResp.Allowed)

	// Fully consistent check (slower)
	consistentResp, _ := client.CheckPermission(ctx, types.CheckRequest{
		Subject:     "user:alice",
		Resource:    "account:ACC-123",
		Permission:  "view",
		Consistency: "fully_consistent",
	})
	log.Printf("Consistent check result: %v", consistentResp.Allowed)

	// Check with snapshot consistency
	if consistentResp.CheckedAt != "" {
		snapshotResp, _ := client.CheckPermission(ctx, types.CheckRequest{
			Subject:     "user:alice",
			Resource:    "account:ACC-123",
			Permission:  "view",
			Consistency: "at_exact_snapshot",
			ZedToken:    consistentResp.CheckedAt,
		})
		log.Printf("Snapshot check result: %v", snapshotResp.Allowed)
	}
}

// BatchExample demonstrates checking multiple permissions
func BatchExample() {
	client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Define multiple checks
	checks := []types.CheckRequest{
		{
			Subject:    "user:alice",
			Resource:   "document:doc-1",
			Permission: "view",
		},
		{
			Subject:    "user:alice",
			Resource:   "document:doc-2",
			Permission: "edit",
		},
		{
			Subject:    "user:bob",
			Resource:   "document:doc-1",
			Permission: "delete",
		},
	}

	// Execute all checks
	results := make([]types.CheckResponse, len(checks))
	for i, check := range checks {
		resp, err := client.CheckPermission(ctx, check)
		if err != nil {
			log.Printf("Error checking %s on %s: %v", check.Subject, check.Resource, err)
			continue
		}
		results[i] = resp
		log.Printf("%s %s %s: %v", check.Subject, check.Permission, check.Resource, resp.Allowed)
	}
}

// ErrorHandlingExample demonstrates error handling
func ErrorHandlingExample() {
	client, err := spicedb.NewClientWithDefaults("localhost:50051", "devkey")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Invalid request (validation error)
	_, err = client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "invalid-subject", // Missing colon separator
		Resource:   "document:doc-1",
		Permission: "view",
	})
	if spicedb.IsValidationError(err) {
		log.Printf("Validation error caught: %v", err)
	}

	// Example 2: Handle operation errors
	response, err := client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:doc-1",
		Permission: "view",
	})
	if err != nil {
		if spicedb.IsOperationError(err) {
			log.Printf("Operation error: %v", err)
		} else {
			log.Printf("Other error: %v", err)
		}
	} else {
		log.Printf("Permission check succeeded: %v", response.Allowed)
	}
}

// CachingExample demonstrates caching configuration
func CachingExample() {
	// Example 1: Default in-memory cache enabled
	clientWithCache, err := spicedb.NewClient(spicedb.ClientOptions{
		Address:      "localhost:50051",
		PreSharedKey: "devkey",
		// Default: in-memory cache enabled
	})
	if err != nil {
		log.Fatal(err)
	}
	defer clientWithCache.Close()

	// Example 2: Disable caching
	clientNoCache, err := spicedb.NewClient(spicedb.ClientOptions{
		Address:      "localhost:50051",
		PreSharedKey: "devkey",
		DisableCache: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer clientNoCache.Close()

	ctx := context.Background()
	req := types.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:doc-1",
		Permission: "view",
	}

	// First call (cache miss)
	start := time.Now()
	resp1, _ := clientWithCache.CheckPermission(ctx, req)
	duration1 := time.Since(start)
	log.Printf("First call took: %v", duration1)
	log.Printf("Result: %v", resp1.Allowed)

	// Second call (cache hit - should be faster)
	start = time.Now()
	resp2, _ := clientWithCache.CheckPermission(ctx, req)
	duration2 := time.Since(start)
	log.Printf("Second call took: %v (cached)", duration2)
	log.Printf("Result: %v", resp2.Allowed)

	if duration2 < duration1 {
		log.Printf("Cache working: second call was faster!")
	}
}

// TimeoutExample demonstrates context-based timeouts
func TimeoutExample() {
	client, err := spicedb.NewClient(spicedb.ClientOptions{
		Address:        "localhost:50051",
		PreSharedKey:   "devkey",
		RequestTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Example 1: Use default timeout
	ctx := context.Background()
	resp, err := client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:doc-1",
		Permission: "view",
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Result: %v", resp.Allowed)
	}

	// Example 2: Override with custom timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err = client.CheckPermission(ctx, types.CheckRequest{
		Subject:    "user:bob",
		Resource:   "document:doc-2",
		Permission: "edit",
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Result: %v", resp.Allowed)
	}
}

func main() {
	log.Println("Running examples...")
	log.Println("\n=== Basic Example ===")
	BasicExample()

	log.Println("\n=== Builder Example ===")
	BuilderExample()

	log.Println("\n=== Context/Caveat Example ===")
	ContextExample()

	log.Println("\n=== Consistency Example ===")
	ConsistencyExample()

	log.Println("\n=== Batch Example ===")
	BatchExample()

	log.Println("\n=== Error Handling Example ===")
	ErrorHandlingExample()

	log.Println("\n=== Caching Example ===")
	CachingExample()

	log.Println("\n=== Timeout Example ===")
	TimeoutExample()
}
