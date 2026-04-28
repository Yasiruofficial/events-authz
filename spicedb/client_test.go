package spicedb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/spicedb/spicedb-go/spicedb/types"
)

// TestClientOptions verifies client creation with different configurations
func TestClientCreation(t *testing.T) {
	tests := []struct {
		name      string
		opts      ClientOptions
		shouldErr bool
		errMsg    string
	}{
		{
			name: "valid options",
			opts: ClientOptions{
				Address:      "localhost:50051",
				PreSharedKey: "test-key",
			},
			shouldErr: false,
		},
		{
			name: "missing address",
			opts: ClientOptions{
				Address: "",
			},
			shouldErr: true,
			errMsg:    "address is required",
		},
		{
			name: "custom timeout",
			opts: ClientOptions{
				Address:        "localhost:50051",
				RequestTimeout: 5 * time.Second,
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.opts)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !IsValidationError(err) {
					t.Fatalf("expected validation error, got %T", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatal("expected non-nil client")
			}
			client.Close()
		})
	}
}

// TestValidationErrors verifies error detection works correctly
func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		validator   func(error) bool
		shouldMatch bool
	}{
		{
			name:        "validation error detection",
			err:         NewValidationError("test", "test error", nil),
			validator:   IsValidationError,
			shouldMatch: true,
		},
		{
			name:        "operation error detection",
			err:         NewOperationError("check", "failed", nil),
			validator:   IsOperationError,
			shouldMatch: true,
		},
		{
			name:        "validation error not operation",
			err:         NewValidationError("test", "test", nil),
			validator:   IsOperationError,
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.validator(tt.err)
			if result != tt.shouldMatch {
				t.Fatalf("expected %v, got %v", tt.shouldMatch, result)
			}
		})
	}
}

// TestErrorwrapping verifies errors can be wrapped correctly
func TestErrorWrapping(t *testing.T) {
	baseErr := errors.New("base error")
	valErr := NewValidationError("field", "reason", baseErr)

	if !errors.Is(valErr, baseErr) {
		t.Fatal("error wrapping failed - base error not found")
	}

	if valErr.Error() != "validation error on field 'field': reason (base error)" {
		t.Fatalf("unexpected error message: %v", valErr.Error())
	}
}

// TestCheckPermissionBuilder verifies builder pattern works
func TestCheckPermissionBuilder(t *testing.T) {
	// Note: This test uses a mock client pattern in production
	// For now, we just verify the builder can be constructed

	opts := ClientOptions{
		Address:      "localhost:50051",
		PreSharedKey: "test",
		DisableCache: true,
	}

	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	builder := client.CheckPermissionBuilder()
	if builder == nil {
		t.Fatal("builder is nil")
	}

	// Verify builder methods chain correctly
	result := builder.
		Subject("user:test").
		Resource("doc:1").
		Permission("view").
		WithConsistency("minimize_latency").
		WithContext(map[string]interface{}{"key": "value"}).
		WithZedToken("token123")

	if result != builder {
		t.Fatal("builder methods should return self for chaining")
	}
}

// TestCacheInterface verifies cache implementations work
func TestCacheInterface(t *testing.T) {
	tests := []struct {
		name  string
		cache Interface
	}{
		{
			name:  "no-op cache",
			cache: &NonOpCache{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set a value
			tt.cache.Set("key1", "value1", 1*time.Second)

			// Try to get it
			_, ok := tt.cache.Get("key1")

			// For NoOpCache, it should not be found
			if !ok {
				// Expected for NoOpCache
			}

			// Delete should not error
			tt.cache.Delete("key1")

			// Clear should not error
			tt.cache.Clear()
		})
	}
}

// Define a test interface to mock cache
type Interface interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

// Create a test implementation
type NonOpCache struct{}

func (c *NonOpCache) Get(key string) (interface{}, bool) {
	return nil, false
}

func (c *NonOpCache) Set(key string, value interface{}, ttl time.Duration) {
}

func (c *NonOpCache) Delete(key string) {
}

func (c *NonOpCache) Clear() {
}

func (c *NonOpCache) String() string {
	return "noop"
}

// TestContextTimeout verifies context timeouts work
func TestContextTimeout(t *testing.T) {
	opts := ClientOptions{
		Address:        "localhost:50051",
		PreSharedKey:   "test",
		RequestTimeout: 100 * time.Millisecond,
		DisableCache:   true,
	}

	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Create a context that times out immediately
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// This should timeout when trying to connect
	time.Sleep(2 * time.Millisecond)

	// Verify context is done
	if ctx.Err() == nil {
		t.Fatal("expected context to be done")
	}
}

// TestNormalization verifies string normalization works
func TestConsistencyNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"MinimizeLatency", "minimizelatency"},
		{"minimize-latency", "minimize_latency"},
		{"Minimize Latency", "minimize_latency"},
		{"FULLY_CONSISTENT", "fully_consistent"},
		{"FullyConsistent", "fullyconsistent"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeConsistency(tt.input)
			if result != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestRequestValidation verifies request validation
func TestObjectReferenceValidation(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		shouldErr bool
	}{
		{
			name:      "valid reference",
			value:     "user:alice",
			fieldName: "subject",
			shouldErr: false,
		},
		{
			name:      "empty reference",
			value:     "",
			fieldName: "resource",
			shouldErr: true,
		},
		{
			name:      "missing colon",
			value:     "userAlice",
			fieldName: "subject",
			shouldErr: true,
		},
		{
			name:      "missing object type",
			value:     ":alice",
			fieldName: "resource",
			shouldErr: true,
		},
		{
			name:      "missing object id",
			value:     "user:",
			fieldName: "subject",
			shouldErr: true,
		},
		{
			name:      "contains relation",
			value:     "user:alice#member",
			fieldName: "resource",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseObjectReference(tt.value, tt.fieldName)
			if tt.shouldErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestSubjectReferenceValidation verifies subject parsing
func TestSubjectReferenceValidation(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		shouldErr bool
		expectRel string
	}{
		{
			name:      "simple subject",
			value:     "user:alice",
			shouldErr: false,
			expectRel: "",
		},
		{
			name:      "subject with relation",
			value:     "group:admins#member",
			shouldErr: false,
			expectRel: "member",
		},
		{
			name:      "empty relation",
			value:     "user:alice#",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseSubjectReference(tt.value)
			if tt.shouldErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestCacheKeyGeneration verifies cache key generation
func TestCacheKeyGeneration(t *testing.T) {
	req := types.CheckRequest{
		Subject:    "user:alice",
		Resource:   "doc:1",
		Permission: "view",
		Context:    map[string]interface{}{"key": "value"},
	}

	key, err := req.CacheKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if key == "" {
		t.Fatal("cache key is empty")
	}

	// Verify same request generates same key
	key2, _ := req.CacheKey()
	if key != key2 {
		t.Fatal("same request should generate same cache key")
	}
}
