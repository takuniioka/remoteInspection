package domain

import "fmt"

// Package domain defines core domain types and errors used across the
// backend. This file contains an `AppError` type which standardizes error
// reporting by attaching a machine-friendly error code, human message,
// optional details map, and an underlying error for debugging.

// AppError represents application errors with a stable error code, a
// human-readable message, optional structured details, and an underlying
// wrapped error. Use this type for errors that will be returned from
// service/repository layers so HTTP handlers can map them to appropriate
// HTTP status codes and response bodies.
type AppError struct {
	Code    string
	Message string
	Details map[string]interface{}
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common error codes used by `AppError.Code`. These are intentionally
// simple string constants so they can be serialized over APIs and used by
// clients to implement error handling logic.
const (
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeInvalidRequest  = "INVALID_REQUEST"
	ErrCodeConflict        = "CONFLICT"
	ErrCodeInternal        = "INTERNAL_ERROR"
	ErrCodeS3Upload        = "S3_UPLOAD_ERROR"
	ErrCodeDynamoDB        = "DYNAMODB_ERROR"
	ErrCodeCognito         = "COGNITO_ERROR"
	ErrCodeTimeout         = "TIMEOUT"
	ErrCodeVideoProvider   = "VIDEO_PROVIDER_ERROR"
)

// NewAppError creates a new AppError with the provided code, message and
// optional wrapped error. The Details map is initialized empty and can be
// enriched later via `WithDetails` to include structured context useful
// for logging or for clients.
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}

// WithDetails attaches a details map to the error and returns the same
// `AppError` instance to allow fluent usage. `details` should contain
// serializable key/value pairs that explain contextual information
// (e.g., request id, resource id) helpful during debugging.
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}
