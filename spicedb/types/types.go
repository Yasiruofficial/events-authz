package types

import (
	"encoding/json"
)

// CheckRequest represents a permission check request
type CheckRequest struct {
	// Subject is the entity checking permission (e.g., "user:alice" or "user:alice#member")
	Subject string `json:"subject"`
	// Resource is the object being accessed (e.g., "document:budget-2026")
	Resource string `json:"resource"`
	// Permission is the permission being checked (e.g., "view", "edit")
	Permission string `json:"permission"`
	// Context is optional caveat context for the check
	Context map[string]interface{} `json:"context,omitempty"`
	// Consistency specifies consistency guarantees (minimize_latency, fully_consistent, at_least_as_fresh, at_exact_snapshot)
	Consistency string `json:"consistency,omitempty"`
	// ZedToken is required for at_least_as_fresh and at_exact_snapshot consistency modes
	ZedToken string `json:"zed_token,omitempty"`
}

// CheckResponse represents the response from a permission check
type CheckResponse struct {
	// Allowed indicates if the permission is granted
	Allowed bool `json:"allowed"`
	// Permissionship provides detailed permission status (has_permission, no_permission, conditional_permission)
	Permissionship string `json:"permissionship,omitempty"`
	// CheckedAt contains the ZedToken for consistency tracking
	CheckedAt string `json:"checked_at,omitempty"`
}

// CacheKey generates a cache key from the check request
func (r CheckRequest) CacheKey() (string, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

// RelationshipFilter represents filtering criteria for relationships
type RelationshipFilter struct {
	ResourceType       string
	ResourceID         string
	Relation           string
	SubjectType        string
	SubjectID          string
	SubjectRelation    string
	OptionalAttributes map[string]interface{}
}

// Relationship represents a relationship in SpiceDB
type Relationship struct {
	Resource             string
	Relation             string
	Subject              string
	OptionalCaveatName   string
	OptionalCaveatFields map[string]interface{}
}

// LookupResourcesRequest represents a lookup resources request
type LookupResourcesRequest struct {
	ResourceType      string
	Permission        string
	Subject           string
	Context           map[string]interface{}
	Consistency       string
	ZedToken          string
	OptionalBatchSize int32
	OptionalLimit     int32
}

// LookupResourcesResponse represents a response from lookup resources
type LookupResourcesResponse struct {
	ResourceID     string
	Permissionship string
	CheckedAt      string
}

// LookupSubjectsRequest represents a lookup subjects request
type LookupSubjectsRequest struct {
	ResourceType      string
	ResourceID        string
	Permission        string
	SubjectType       string
	Context           map[string]interface{}
	Consistency       string
	ZedToken          string
	OptionalBatchSize int32
	OptionalLimit     int32
}

// LookupSubjectsResponse represents a response from lookup subjects
type LookupSubjectsResponse struct {
	SubjectID          string
	Permissionship     string
	PartialCaveatMatch bool
	CheckedAt          string
}

// ReadRelationshipsRequest represents a request to read relationships
type ReadRelationshipsRequest struct {
	Filter        RelationshipFilter
	Consistency   string
	ZedToken      string
	OptionalLimit int32
}

// WriteRelationshipsRequest represents a request to write (create/delete) relationships
type WriteRelationshipsRequest struct {
	Updates []*Relationship
}

// DeleteRelationshipsRequest represents a request to delete relationships
type DeleteRelationshipsRequest struct {
	Filter                RelationshipFilter
	OptionalPreconditions []*Relationship
}

// ExpireRelationshipsRequest represents a request to set relationship expiration
type ExpireRelationshipsRequest struct {
	Filter    RelationshipFilter
	ExpiresAt string
}
