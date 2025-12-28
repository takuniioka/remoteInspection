package domain

import "fmt"

// AppError represents application errors with code and message
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

// Common error codes
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

// NewAppError creates a new app error
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to an app error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}
