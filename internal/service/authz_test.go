package service

import (
	"context"
	"testing"

	"events-authz/internal/cache"
	"events-authz/internal/model"
)

type fakePermissionChecker struct {
	calls    int
	response model.CheckResponse
	err      error
}

func (f *fakePermissionChecker) CheckPermission(_ context.Context, _ model.CheckRequest) (model.CheckResponse, error) {
	f.calls++
	return f.response, f.err
}

func TestAuthzServiceCachesIdenticalRequests(t *testing.T) {
	checker := &fakePermissionChecker{
		response: model.CheckResponse{Allowed: true, Permissionship: "has_permission"},
	}
	service := NewAuthzService(checker, cache.New())
	request := model.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:budget-2026",
		Permission: "view",
	}

	first, err := service.Check(context.Background(), request)
	if err != nil {
		t.Fatalf("first Check returned error: %v", err)
	}
	second, err := service.Check(context.Background(), request)
	if err != nil {
		t.Fatalf("second Check returned error: %v", err)
	}

	if checker.calls != 1 {
		t.Fatalf("expected exactly one backend call, got %d", checker.calls)
	}
	if first != second {
		t.Fatalf("expected cached response to match first response, got %#v and %#v", first, second)
	}
}

func TestAuthzServiceCacheKeyIncludesContext(t *testing.T) {
	checker := &fakePermissionChecker{
		response: model.CheckResponse{Allowed: true, Permissionship: "has_permission"},
	}
	service := NewAuthzService(checker, cache.New())

	_, err := service.Check(context.Background(), model.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:budget-2026",
		Permission: "view",
		Context:    map[string]any{"region": "eu"},
	})
	if err != nil {
		t.Fatalf("first Check returned error: %v", err)
	}

	_, err = service.Check(context.Background(), model.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:budget-2026",
		Permission: "view",
		Context:    map[string]any{"region": "us"},
	})
	if err != nil {
		t.Fatalf("second Check returned error: %v", err)
	}

	if checker.calls != 2 {
		t.Fatalf("expected context changes to bypass cache, got %d backend calls", checker.calls)
	}
}
