package spicedb

import (
	"context"

	"github.com/Yasiruofficial/events-authz/spicedb/types"
)

// CheckPermissionBuilder provides a fluent interface for building permission checks
type CheckPermissionBuilder struct {
	client      *Client
	subject     string
	resource    string
	permission  string
	context     map[string]interface{}
	consistency string
	zedToken    string
}

// Subject sets the subject for the permission check
func (b *CheckPermissionBuilder) Subject(subject string) *CheckPermissionBuilder {
	b.subject = subject
	return b
}

// Resource sets the resource for the permission check
func (b *CheckPermissionBuilder) Resource(resource string) *CheckPermissionBuilder {
	b.resource = resource
	return b
}

// Permission sets the permission for the check
func (b *CheckPermissionBuilder) Permission(permission string) *CheckPermissionBuilder {
	b.permission = permission
	return b
}

// WithContext sets the context (caveat fields) for the check
func (b *CheckPermissionBuilder) WithContext(ctx map[string]interface{}) *CheckPermissionBuilder {
	b.context = ctx
	return b
}

// WithConsistency sets the consistency level for the check
// Valid values: "minimize_latency", "fully_consistent", "at_least_as_fresh", "at_exact_snapshot"
func (b *CheckPermissionBuilder) WithConsistency(consistency string) *CheckPermissionBuilder {
	b.consistency = consistency
	return b
}

// WithZedToken sets the zed token for snapshot-based consistency
func (b *CheckPermissionBuilder) WithZedToken(token string) *CheckPermissionBuilder {
	b.zedToken = token
	return b
}

// Check executes the permission check
func (b *CheckPermissionBuilder) Check(ctx context.Context) (types.CheckResponse, error) {
	req := types.CheckRequest{
		Subject:     b.subject,
		Resource:    b.resource,
		Permission:  b.permission,
		Context:     b.context,
		Consistency: b.consistency,
		ZedToken:    b.zedToken,
	}
	return b.client.CheckPermission(ctx, req)
}

// IsAllowed is a convenience method that checks permission and returns only the allowed status
func (b *CheckPermissionBuilder) IsAllowed(ctx context.Context) (bool, error) {
	resp, err := b.Check(ctx)
	if err != nil {
		return false, err
	}
	return resp.Allowed, nil
}
