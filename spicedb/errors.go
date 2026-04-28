package spicedb

import (
	"errors"
	"fmt"
)

// Error types for SpiceDB operations
var (
	// ErrInvalidRequest is returned when request validation fails
	ErrInvalidRequest = errors.New("invalid request")
	// ErrPermissionDenied is returned when access is denied
	ErrPermissionDenied = errors.New("permission denied")
	// ErrRelationshipNotFound is returned when a relationship is not found
	ErrRelationshipNotFound = errors.New("relationship not found")
	// ErrInconsistentRead is returned when consistency check fails
	ErrInconsistentRead = errors.New("inconsistent read")
	// ErrConnectionFailed is returned when connection to SpiceDB fails
	ErrConnectionFailed = errors.New("connection failed")
	// ErrDeadlineExceeded is returned when the request times out
	ErrDeadlineExceeded = errors.New("deadline exceeded")
)

// ValidationError represents a validation error with field information
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation error on field '%s': %s (%v)", e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, err error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Err:     err,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// OperationError represents an error from a SpiceDB operation
type OperationError struct {
	Operation string
	Message   string
	Err       error
}

func (e *OperationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("operation '%s' failed: %s (%v)", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("operation '%s' failed: %s", e.Operation, e.Message)
}

func (e *OperationError) Unwrap() error {
	return e.Err
}

// NewOperationError creates a new operation error
func NewOperationError(operation, message string, err error) *OperationError {
	return &OperationError{
		Operation: operation,
		Message:   message,
		Err:       err,
	}
}

// IsOperationError checks if an error is an operation error
func IsOperationError(err error) bool {
	var oe *OperationError
	return errors.As(err, &oe)
}
