package spicedb

import (
	"context"
	"errors"
	"testing"
	"time"

	"events-authz/internal/model"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"google.golang.org/grpc"
)

type fakePermissionsServiceClient struct {
	request  *v1.CheckPermissionRequest
	response *v1.CheckPermissionResponse
	err      error
}

func (f *fakePermissionsServiceClient) CheckPermission(_ context.Context, in *v1.CheckPermissionRequest, _ ...grpc.CallOption) (*v1.CheckPermissionResponse, error) {
	f.request = in
	return f.response, f.err
}

func TestClientCheckPermissionBuildsSpiceDBRequest(t *testing.T) {
	fake := &fakePermissionsServiceClient{
		response: &v1.CheckPermissionResponse{
			Permissionship: v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION,
			CheckedAt:      &v1.ZedToken{Token: "zed-token"},
		},
	}

	client := newClient(fake, nil, time.Second, "fully_consistent")
	response, err := client.CheckPermission(context.Background(), model.CheckRequest{
		Subject:    "user:alice",
		Resource:   "document:budget-2026",
		Permission: "view",
		Context: map[string]any{
			"ip": "127.0.0.1",
		},
	})
	if err != nil {
		t.Fatalf("CheckPermission returned error: %v", err)
	}

	if !response.Allowed {
		t.Fatalf("expected Allowed to be true")
	}
	if response.Permissionship != "has_permission" {
		t.Fatalf("expected permissionship to be has_permission, got %q", response.Permissionship)
	}
	if response.CheckedAt != "zed-token" {
		t.Fatalf("expected checked_at token to be preserved, got %q", response.CheckedAt)
	}

	if fake.request == nil {
		t.Fatal("expected SpiceDB request to be captured")
	}
	if fake.request.GetResource().GetObjectType() != "document" {
		t.Fatalf("expected resource type document, got %q", fake.request.GetResource().GetObjectType())
	}
	if fake.request.GetResource().GetObjectId() != "budget-2026" {
		t.Fatalf("expected resource id budget-2026, got %q", fake.request.GetResource().GetObjectId())
	}
	if fake.request.GetSubject().GetObject().GetObjectType() != "user" {
		t.Fatalf("expected subject type user, got %q", fake.request.GetSubject().GetObject().GetObjectType())
	}
	if fake.request.GetSubject().GetObject().GetObjectId() != "alice" {
		t.Fatalf("expected subject id alice, got %q", fake.request.GetSubject().GetObject().GetObjectId())
	}
	if !fake.request.GetConsistency().GetFullyConsistent() {
		t.Fatal("expected fully consistent checks by default")
	}
	if got := fake.request.GetContext().GetFields()["ip"].GetStringValue(); got != "127.0.0.1" {
		t.Fatalf("expected caveat context ip to round-trip, got %q", got)
	}
}

func TestClientCheckPermissionRejectsInvalidRequest(t *testing.T) {
	client := newClient(&fakePermissionsServiceClient{}, nil, time.Second, "minimize_latency")

	_, err := client.CheckPermission(context.Background(), model.CheckRequest{
		Subject:    "alice",
		Resource:   "document:budget-2026",
		Permission: "view",
	})
	if !errors.Is(err, ErrInvalidCheckRequest) {
		t.Fatalf("expected ErrInvalidCheckRequest, got %v", err)
	}
}
