package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// AppError represents an application error
type AppError struct {
	Code    int    // HTTP status code
	Message string // Error message
	Err     error  // Original error (if any)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error constructors
func NewBadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, message, err)
}

func NewUnauthorized(message string, err error) *AppError {
	return New(http.StatusUnauthorized, message, err)
}

func NewForbidden(message string, err error) *AppError {
	return New(http.StatusForbidden, message, err)
}

func NewNoContent(message string) *AppError {
	return New(http.StatusNoContent, message, nil)
}

func NewInternalServerError(message string, err error) *AppError {
	return New(http.StatusInternalServerError, message, err)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: message,
	}
}

// HandleError handles an error and returns an appropriate response
func HandleError(w http.ResponseWriter, err error) {
	appErr, ok := IsAppError(err)
	if !ok {
		// If it's not an AppError, wrap it as an internal server error
		appErr = NewInternalServerError("Internal server error", err)
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)

	// Create error response
	response := NewErrorResponse(appErr.Message)
	if err:= json.NewEncoder(w).Encode(response); err != nil {
		logrus.Errorf("failed to write error response %v", err)
	}
}