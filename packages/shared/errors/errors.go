package errors

import (
	"fmt"
)

// Error codes
const (
	ErrInternal         = "INTERNAL_ERROR"
	ErrInternalServer   = "INTERNAL_SERVER_ERROR"
	ErrInvalidInput     = "INVALID_INPUT"
	ErrNotFound         = "NOT_FOUND"
	ErrUnauthorized     = "UNAUTHORIZED"
	ErrForbidden        = "FORBIDDEN"
	ErrConflict         = "CONFLICT"
	ErrTenantNotFound   = "TENANT_NOT_FOUND"
	ErrTenantSuspended  = "TENANT_SUSPENDED"
	ErrInvalidToken     = "INVALID_TOKEN"
	ErrTokenExpired     = "TOKEN_EXPIRED"
	ErrInsufficientPermissions = "INSUFFICIENT_PERMISSIONS"
	ErrRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
)

// APIError represents a structured API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(code, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewAPIErrorWithDetails creates a new API error with details
func NewAPIErrorWithDetails(code, message string, details interface{}) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error constructors
func Internal(message string) *APIError {
	return NewAPIError(ErrInternal, message)
}

func InvalidInput(message string) *APIError {
	return NewAPIError(ErrInvalidInput, message)
}

func NotFound(message string) *APIError {
	return NewAPIError(ErrNotFound, message)
}

func Unauthorized(message string) *APIError {
	return NewAPIError(ErrUnauthorized, message)
}

func Forbidden(message string) *APIError {
	return NewAPIError(ErrForbidden, message)
}

func Conflict(message string) *APIError {
	return NewAPIError(ErrConflict, message)
}
