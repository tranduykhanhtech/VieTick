package errors

import (
    "fmt"
    "net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
    // ErrorTypeValidation represents validation errors
    ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
    // ErrorTypeAuthentication represents authentication errors
    ErrorTypeAuthentication ErrorType = "AUTHENTICATION_ERROR"
    // ErrorTypeAuthorization represents authorization errors
    ErrorTypeAuthorization ErrorType = "AUTHORIZATION_ERROR"
    // ErrorTypeNotFound represents not found errors
    ErrorTypeNotFound ErrorType = "NOT_FOUND_ERROR"
    // ErrorTypeConflict represents conflict errors
    ErrorTypeConflict ErrorType = "CONFLICT_ERROR"
    // ErrorTypeInternal represents internal server errors
    ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
)

// AppError represents an application error
type AppError struct {
    Code      int       `json:"code"`      // HTTP status code
    Type      ErrorType `json:"type"`      // Type of error
    Message   string    `json:"message"`   // User-friendly message
    Details   string    `json:"details"`   // Detailed error message
    Err       error     `json:"-"`         // Original error
    RequestID string    `json:"request_id"` // Request ID for tracking
}

// Error implements the error interface
func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
    }
    return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// New creates a new AppError
func New(code int, errType ErrorType, message, details string, err error) *AppError {
    return &AppError{
        Code:    code,
        Type:    errType,
        Message: message,
        Details: details,
        Err:     err,
    }
}

// Common error constructors with detailed messages

// ValidationError creates a validation error
func ValidationError(message, details string, err error) *AppError {
    return New(
        http.StatusBadRequest,
        ErrorTypeValidation,
        message,
        details,
        err,
    )
}

// AuthenticationError creates an authentication error
func AuthenticationError(message, details string, err error) *AppError {
    return New(
        http.StatusUnauthorized,
        ErrorTypeAuthentication,
        message,
        details,
        err,
    )
}

// AuthorizationError creates an authorization error
func AuthorizationError(message, details string, err error) *AppError {
    return New(
        http.StatusForbidden,
        ErrorTypeAuthorization,
        message,
        details,
        err,
    )
}

// NotFoundError creates a not found error
func NotFoundError(message, details string, err error) *AppError {
    return New(
        http.StatusNotFound,
        ErrorTypeNotFound,
        message,
        details,
        err,
    )
}

// ConflictError creates a conflict error
func ConflictError(message, details string, err error) *AppError {
    return New(
        http.StatusConflict,
        ErrorTypeConflict,
        message,
        details,
        err,
    )
}

// InternalError creates an internal server error
func InternalError(message, details string, err error) *AppError {
    return New(
        http.StatusInternalServerError,
        ErrorTypeInternal,
        message,
        details,
        err,
    )
}

// Common error messages
const (
    // Authentication errors
    ErrInvalidCredentials = "Invalid email or password"
    ErrTokenExpired      = "Authentication token has expired"
    ErrInvalidToken      = "Invalid authentication token"
    ErrMissingToken      = "Authentication token is required"

    // Authorization errors
    ErrInsufficientPoints = "Insufficient points to perform this action"
    ErrDailyLimitReached = "Daily limit for this action has been reached"
    ErrNotVerified       = "User is not verified"

    // Validation errors
    ErrInvalidEmail     = "Invalid email format"
    ErrInvalidPassword  = "Password must be at least 8 characters long"
    ErrInvalidUsername  = "Username must be between 3 and 20 characters"
    ErrRequiredField    = "This field is required"
    ErrInvalidVote      = "Invalid vote value"
    ErrDuplicateVote   = "You have already voted for this answer"

    // Not found errors
    ErrUserNotFound     = "User not found"
    ErrQuestionNotFound = "Question not found"
    ErrAnswerNotFound   = "Answer not found"

    // Conflict errors
    ErrEmailExists      = "Email already exists"
    ErrUsernameExists   = "Username already exists"
    ErrAlreadyVerified  = "Answer is already verified"

    // Internal errors
    ErrDatabase         = "Database error occurred"
    ErrCache           = "Cache error occurred"
    ErrExternalService = "External service error occurred"
) 