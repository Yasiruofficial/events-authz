package spicedb

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spicedb/spicedb-go/internal/model"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	authzed "github.com/authzed/authzed-go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

var ErrInvalidCheckRequest = errors.New("invalid authorization check request")

type permissionsServiceClient interface {
	CheckPermission(ctx context.Context, in *v1.CheckPermissionRequest, opts ...grpc.CallOption) (*v1.CheckPermissionResponse, error)
}

type Client struct {
	permissions        permissionsServiceClient
	closeFn            func() error
	requestTimeout     time.Duration
	defaultConsistency string
}

func NewClient(addr, token string, allowInsecure bool, requestTimeout time.Duration, defaultConsistency string) (*Client, error) {
	options := []grpc.DialOption{transportCredentials(allowInsecure)}
	if token != "" {
		options = append(options, grpc.WithPerRPCCredentials(bearerTokenCredentials{token: token, insecure: allowInsecure}))
	}

	client, err := authzed.NewClient(addr, options...)
	if err != nil {
		return nil, fmt.Errorf("create authzed client: %w", err)
	}

	return newClient(client, client.Close, requestTimeout, defaultConsistency), nil
}

func newClient(permissions permissionsServiceClient, closeFn func() error, requestTimeout time.Duration, defaultConsistency string) *Client {
	if requestTimeout <= 0 {
		requestTimeout = 3 * time.Second
	}

	if strings.TrimSpace(defaultConsistency) == "" {
		defaultConsistency = "minimize_latency"
	}

	return &Client{
		permissions:        permissions,
		closeFn:            closeFn,
		requestTimeout:     requestTimeout,
		defaultConsistency: defaultConsistency,
	}
}

func (c *Client) Close() error {
	if c.closeFn == nil {
		return nil
	}

	return c.closeFn()
}

func (c *Client) CheckPermission(ctx context.Context, req model.CheckRequest) (model.CheckResponse, error) {
	request, err := c.buildCheckPermissionRequest(req)
	if err != nil {
		return model.CheckResponse{}, err
	}

	if err := request.Validate(); err != nil {
		return model.CheckResponse{}, fmt.Errorf("%w: %v", ErrInvalidCheckRequest, err)
	}

	if c.requestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.requestTimeout)
		defer cancel()
	}

	response, err := c.permissions.CheckPermission(ctx, request)
	if err != nil {
		return model.CheckResponse{}, fmt.Errorf("check permission in spicedb: %w", err)
	}

	result := model.CheckResponse{
		Allowed:        response.GetPermissionship() == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION,
		Permissionship: permissionshipName(response.GetPermissionship()),
	}

	if checkedAt := response.GetCheckedAt(); checkedAt != nil {
		result.CheckedAt = checkedAt.GetToken()
	}

	return result, nil
}

func (c *Client) buildCheckPermissionRequest(req model.CheckRequest) (*v1.CheckPermissionRequest, error) {
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

func buildContextStruct(values map[string]any) (*structpb.Struct, error) {
	if len(values) == 0 {
		return nil, nil
	}

	contextStruct, err := structpb.NewStruct(values)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid caveat context: %v", ErrInvalidCheckRequest, err)
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
			return nil, fmt.Errorf("%w: zed_token is required for at_least_as_fresh consistency", ErrInvalidCheckRequest)
		}
		return &v1.Consistency{Requirement: &v1.Consistency_AtLeastAsFresh{AtLeastAsFresh: &v1.ZedToken{Token: strings.TrimSpace(zedToken)}}}, nil
	case "at_exact_snapshot":
		if strings.TrimSpace(zedToken) == "" {
			return nil, fmt.Errorf("%w: zed_token is required for at_exact_snapshot consistency", ErrInvalidCheckRequest)
		}
		return &v1.Consistency{Requirement: &v1.Consistency_AtExactSnapshot{AtExactSnapshot: &v1.ZedToken{Token: strings.TrimSpace(zedToken)}}}, nil
	default:
		return nil, fmt.Errorf("%w: unsupported consistency %q", ErrInvalidCheckRequest, invalidValue)
	}
}

func parseObjectReference(value, fieldName string) (*v1.ObjectReference, error) {
	reference := strings.TrimSpace(value)
	if reference == "" {
		return nil, fmt.Errorf("%w: %s must be in type:id format", ErrInvalidCheckRequest, fieldName)
	}

	if strings.Contains(reference, "#") {
		return nil, fmt.Errorf("%w: %s must not contain a relation; use type:id format", ErrInvalidCheckRequest, fieldName)
	}

	objectType, objectID, ok := strings.Cut(reference, ":")
	if !ok || strings.TrimSpace(objectType) == "" || strings.TrimSpace(objectID) == "" {
		return nil, fmt.Errorf("%w: %s must be in type:id format", ErrInvalidCheckRequest, fieldName)
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
			return nil, fmt.Errorf("%w: subject relation must not be empty", ErrInvalidCheckRequest)
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

func transportCredentials(allowInsecure bool) grpc.DialOption {
	if allowInsecure {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	return grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12}))
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
