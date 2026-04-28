package spicedb

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/Yasiruofficial/events-authz/spicedb/cache"
	"github.com/Yasiruofficial/events-authz/spicedb/types"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	authzed "github.com/authzed/authzed-go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

// permissionsServiceClient defines the interface for SpiceDB permissions service
type permissionsServiceClient interface {
	CheckPermission(ctx context.Context, in *v1.CheckPermissionRequest, opts ...grpc.CallOption) (*v1.CheckPermissionResponse, error)
	ReadRelationships(ctx context.Context, in *v1.ReadRelationshipsRequest, opts ...grpc.CallOption) (v1.PermissionsService_ReadRelationshipsClient, error)
	WriteRelationships(ctx context.Context, in *v1.WriteRelationshipsRequest, opts ...grpc.CallOption) (*v1.WriteRelationshipsResponse, error)
	DeleteRelationships(ctx context.Context, in *v1.DeleteRelationshipsRequest, opts ...grpc.CallOption) (*v1.DeleteRelationshipsResponse, error)
	LookupResources(ctx context.Context, in *v1.LookupResourcesRequest, opts ...grpc.CallOption) (v1.PermissionsService_LookupResourcesClient, error)
	LookupSubjects(ctx context.Context, in *v1.LookupSubjectsRequest, opts ...grpc.CallOption) (v1.PermissionsService_LookupSubjectsClient, error)
}

// ClientOptions contains configuration for the SpiceDB client
type ClientOptions struct {
	// Address of the SpiceDB server (required)
	Address string
	// PreSharedKey for authentication
	PreSharedKey string
	// TLSEnabled enables TLS for the connection
	TLSEnabled bool
	// InsecureSkipVerify skips TLS verification (only if TLSEnabled is true)
	InsecureSkipVerify bool
	// RequestTimeout is the default timeout for requests
	RequestTimeout time.Duration
	// DefaultConsistency is the default consistency level for requests
	DefaultConsistency string
	// Cache is an optional cache implementation (defaults to in-memory)
	Cache cache.Interface
	// DisableCache disables caching entirely
	DisableCache bool
}

// Client is the main SpiceDB client
type Client struct {
	permissions        permissionsServiceClient
	closeFn            func() error
	requestTimeout     time.Duration
	defaultConsistency string
	cache              cache.Interface
	disableCache       bool
}

// NewClient creates a new SpiceDB client with the given options
func NewClient(opts ClientOptions) (*Client, error) {
	if strings.TrimSpace(opts.Address) == "" {
		return nil, NewValidationError("Address", "address is required", nil)
	}

	// Build gRPC dial options
	dialOpts := []grpc.DialOption{transportCredentials(opts.TLSEnabled, opts.InsecureSkipVerify)}

	if opts.PreSharedKey != "" {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(
			bearerTokenCredentials{
				token:    opts.PreSharedKey,
				insecure: !opts.TLSEnabled,
			},
		))
	}

	// Create authzed client
	client, err := authzed.NewClient(opts.Address, dialOpts...)
	if err != nil {
		return nil, NewOperationError("initialize", fmt.Sprintf("failed to create SpiceDB client at %s", opts.Address), err)
	}

	// Normalize and validate consistency
	defaultConsistency := normalizeConsistency(opts.DefaultConsistency)
	if defaultConsistency == "" {
		defaultConsistency = "minimize_latency"
	}

	// Set up caching
	var c cache.Interface
	if opts.DisableCache {
		c = &cache.NoOpCache{}
	} else if opts.Cache != nil {
		c = opts.Cache
	} else {
		c = cache.NewInMemoryCache()
	}

	// Set request timeout
	requestTimeout := opts.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = 3 * time.Second
	}

	return &Client{
		permissions:        client,
		closeFn:            client.Close,
		requestTimeout:     requestTimeout,
		defaultConsistency: defaultConsistency,
		cache:              c,
		disableCache:       opts.DisableCache,
	}, nil
}

// NewClientWithDefaults creates a new client with sensible defaults for insecure development
func NewClientWithDefaults(address, preSharedKey string) (*Client, error) {
	return NewClient(ClientOptions{
		Address:            address,
		PreSharedKey:       preSharedKey,
		TLSEnabled:         false,
		RequestTimeout:     3 * time.Second,
		DefaultConsistency: "minimize_latency",
	})
}

// Close closes the client connection
func (c *Client) Close() error {
	if c.closeFn == nil {
		return nil
	}
	return c.closeFn()
}

// CheckPermission checks if a subject has a permission on a resource
func (c *Client) CheckPermission(ctx context.Context, req types.CheckRequest) (types.CheckResponse, error) {
	// Check cache first
	if !c.disableCache {
		cacheKey, err := req.CacheKey()
		if err == nil {
			if val, ok := c.cache.Get(cacheKey); ok {
				if cached, ok := val.(types.CheckResponse); ok {
					return cached, nil
				}
			}
		}
	}

	// Build the request
	grpcReq, err := c.buildCheckPermissionRequest(req)
	if err != nil {
		return types.CheckResponse{}, err
	}

	// Validate the built request
	if err := grpcReq.Validate(); err != nil {
		return types.CheckResponse{}, NewValidationError("CheckPermissionRequest", "request validation failed", err)
	}

	// Apply timeout
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	// Execute the check
	response, err := c.permissions.CheckPermission(ctx, grpcReq)
	if err != nil {
		return types.CheckResponse{}, NewOperationError("CheckPermission", "failed to check permission", err)
	}

	// Build response
	result := types.CheckResponse{
		Allowed:        response.GetPermissionship() == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION,
		Permissionship: permissionshipName(response.GetPermissionship()),
	}

	if checkedAt := response.GetCheckedAt(); checkedAt != nil {
		result.CheckedAt = checkedAt.GetToken()
	}

	// Cache the result
	if !c.disableCache {
		cacheKey, err := req.CacheKey()
		if err == nil && cacheKey != "" {
			c.cache.Set(cacheKey, result, 5*time.Second)
		}
	}

	return result, nil
}

// CheckPermissionBuilder creates a builder for permission checks
func (c *Client) CheckPermissionBuilder() *CheckPermissionBuilder {
	return &CheckPermissionBuilder{
		client:      c,
		consistency: c.defaultConsistency,
	}
}

// buildCheckPermissionRequest converts a CheckRequest to a gRPC CheckPermissionRequest
func (c *Client) buildCheckPermissionRequest(req types.CheckRequest) (*v1.CheckPermissionRequest, error) {
	resource, err := parseObjectReference(req.Resource, "resource")
	if err != nil {
		return nil, err
	}

	subject, err := parseSubjectReference(req.Subject)
	if err != nil {
		return nil, err
	}

	consistency, err := buildConsistency(req.Consistency, c.defaultConsistency, req.ZedToken)
	if err != nil {
		return nil, err
	}

	contextStruct, err := buildContextStruct(req.Context)
	if err != nil {
		return nil, err
	}

	return &v1.CheckPermissionRequest{
		Resource:    resource,
		Permission:  strings.TrimSpace(req.Permission),
		Subject:     subject,
		Context:     contextStruct,
		Consistency: consistency,
	}, nil
}

// Helper functions for building requests

func buildContextStruct(values map[string]interface{}) (*structpb.Struct, error) {
	if len(values) == 0 {
		return nil, nil
	}

	contextStruct, err := structpb.NewStruct(values)
	if err != nil {
		return nil, NewValidationError("Context", "invalid caveat context", err)
	}

	return contextStruct, nil
}

func buildConsistency(requestConsistency, defaultConsistency, zedToken string) (*v1.Consistency, error) {
	consistency := normalizeConsistency(requestConsistency)
	invalidValue := requestConsistency
	if consistency == "" {
		consistency = normalizeConsistency(defaultConsistency)
		invalidValue = defaultConsistency
	}

	switch consistency {
	case "", "minimize_latency":
		return &v1.Consistency{Requirement: &v1.Consistency_MinimizeLatency{MinimizeLatency: true}}, nil
	case "fully_consistent":
		return &v1.Consistency{Requirement: &v1.Consistency_FullyConsistent{FullyConsistent: true}}, nil
	case "at_least_as_fresh":
		if strings.TrimSpace(zedToken) == "" {
			return nil, NewValidationError("ZedToken", "zed_token is required for at_least_as_fresh consistency", nil)
		}
		return &v1.Consistency{Requirement: &v1.Consistency_AtLeastAsFresh{AtLeastAsFresh: &v1.ZedToken{Token: strings.TrimSpace(zedToken)}}}, nil
	case "at_exact_snapshot":
		if strings.TrimSpace(zedToken) == "" {
			return nil, NewValidationError("ZedToken", "zed_token is required for at_exact_snapshot consistency", nil)
		}
		return &v1.Consistency{Requirement: &v1.Consistency_AtExactSnapshot{AtExactSnapshot: &v1.ZedToken{Token: strings.TrimSpace(zedToken)}}}, nil
	default:
		return nil, NewValidationError("Consistency", fmt.Sprintf("unsupported consistency %q", invalidValue), nil)
	}
}

func parseObjectReference(value, fieldName string) (*v1.ObjectReference, error) {
	reference := strings.TrimSpace(value)
	if reference == "" {
		return nil, NewValidationError(fieldName, "must be in type:id format", nil)
	}

	if strings.Contains(reference, "#") {
		return nil, NewValidationError(fieldName, "must not contain a relation; use type:id format", nil)
	}

	objectType, objectID, ok := strings.Cut(reference, ":")
	if !ok || strings.TrimSpace(objectType) == "" || strings.TrimSpace(objectID) == "" {
		return nil, NewValidationError(fieldName, "must be in type:id format", nil)
	}

	return &v1.ObjectReference{
		ObjectType: strings.TrimSpace(objectType),
		ObjectId:   strings.TrimSpace(objectID),
	}, nil
}

func parseSubjectReference(value string) (*v1.SubjectReference, error) {
	reference := strings.TrimSpace(value)
	baseReference, relation, hasRelation := strings.Cut(reference, "#")

	object, err := parseObjectReference(baseReference, "subject")
	if err != nil {
		return nil, err
	}

	if hasRelation {
		relation = strings.TrimSpace(relation)
		if relation == "" {
			return nil, NewValidationError("subject", "relation must not be empty", nil)
		}
	}

	return &v1.SubjectReference{
		Object:           object,
		OptionalRelation: relation,
	}, nil
}

func normalizeConsistency(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	return normalized
}

func permissionshipName(value v1.CheckPermissionResponse_Permissionship) string {
	return strings.ToLower(strings.TrimPrefix(value.String(), "PERMISSIONSHIP_"))
}

func transportCredentials(tlsEnabled bool, insecureSkipVerify bool) grpc.DialOption {
	if !tlsEnabled {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: insecureSkipVerify,
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
}

type bearerTokenCredentials struct {
	token    string
	insecure bool
}

func (c bearerTokenCredentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer " + c.token}, nil
}

func (c bearerTokenCredentials) RequireTransportSecurity() bool {
	return !c.insecure
}
