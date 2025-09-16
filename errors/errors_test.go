package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	t.Run("should return message when no wrapped error", func(t *testing.T) {
		appErr := New(http.StatusBadRequest, "Bad request", nil)
		assert.Equal(t, "Bad request", appErr.Error())
	})

	t.Run("should return message with wrapped error", func(t *testing.T) {
		originalErr := errors.New("original error")
		appErr := New(http.StatusBadRequest, "Bad request", originalErr)
		assert.Equal(t, "Bad request: original error", appErr.Error())
	})
}

func TestAppError_Unwrap(t *testing.T) {
	t.Run("should return wrapped error", func(t *testing.T) {
		originalErr := errors.New("original error")
		appErr := New(http.StatusBadRequest, "Bad request", originalErr)
		assert.Equal(t, originalErr, appErr.Unwrap())
	})

	t.Run("should return nil when no wrapped error", func(t *testing.T) {
		appErr := New(http.StatusBadRequest, "Bad request", nil)
		assert.Nil(t, appErr.Unwrap())
	})
}

func TestNew(t *testing.T) {
	t.Run("should create AppError with all fields", func(t *testing.T) {
		originalErr := errors.New("test error")
		appErr := New(http.StatusInternalServerError, "Test message", originalErr)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, "Test message", appErr.Message)
		assert.Equal(t, originalErr, appErr.Err)
	})

	t.Run("should create AppError without wrapped error", func(t *testing.T) {
		appErr := New(http.StatusOK, "Success", nil)
		assert.Equal(t, http.StatusOK, appErr.Code)
		assert.Equal(t, "Success", appErr.Message)
		assert.Nil(t, appErr.Err)
	})
}

func TestNewBadRequest(t *testing.T) {
	t.Run("should create bad request error", func(t *testing.T) {
		originalErr := errors.New("validation error")
		appErr := NewBadRequest("Invalid input", originalErr)
		assert.Equal(t, http.StatusBadRequest, appErr.Code)
		assert.Equal(t, "Invalid input", appErr.Message)
		assert.Equal(t, originalErr, appErr.Err)
	})
}

func TestNewUnauthorized(t *testing.T) {
	t.Run("should create unauthorized error", func(t *testing.T) {
		originalErr := errors.New("auth error")
		appErr := NewUnauthorized("Authentication required", originalErr)
		assert.Equal(t, http.StatusUnauthorized, appErr.Code)
		assert.Equal(t, "Authentication required", appErr.Message)
		assert.Equal(t, originalErr, appErr.Err)
	})
}

func TestNewForbidden(t *testing.T) {
	t.Run("should create forbidden error", func(t *testing.T) {
		originalErr := errors.New("permission error")
		appErr := NewForbidden("Access denied", originalErr)
		assert.Equal(t, http.StatusForbidden, appErr.Code)
		assert.Equal(t, "Access denied", appErr.Message)
		assert.Equal(t, originalErr, appErr.Err)
	})
}

func TestNewNoContent(t *testing.T) {
	t.Run("should create no content error", func(t *testing.T) {
		appErr := NewNoContent("No content available")
		assert.Equal(t, http.StatusNoContent, appErr.Code)
		assert.Equal(t, "No content available", appErr.Message)
		assert.Nil(t, appErr.Err)
	})
}

func TestNewInternalServerError(t *testing.T) {
	t.Run("should create internal server error", func(t *testing.T) {
		originalErr := errors.New("database error")
		appErr := NewInternalServerError("Something went wrong", originalErr)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		assert.Equal(t, "Something went wrong", appErr.Message)
		assert.Equal(t, originalErr, appErr.Err)
	})
}

func TestIsAppError(t *testing.T) {
	t.Run("should return true for AppError", func(t *testing.T) {
		appErr := New(http.StatusBadRequest, "Test", nil)
		result, ok := IsAppError(appErr)
		assert.True(t, ok)
		assert.Equal(t, appErr, result)
	})

	t.Run("should return false for regular error", func(t *testing.T) {
		err := errors.New("regular error")
		result, ok := IsAppError(err)
		assert.False(t, ok)
		assert.Nil(t, result)
	})

	t.Run("should return false for nil error", func(t *testing.T) {
		result, ok := IsAppError(nil)
		assert.False(t, ok)
		assert.Nil(t, result)
	})
}

func TestNewErrorResponse(t *testing.T) {
	t.Run("should create error response", func(t *testing.T) {
		response := NewErrorResponse("Test error message")
		assert.Equal(t, "Test error message", response.Error)
	})
}

func TestHandleError(t *testing.T) {
	t.Run("should handle AppError correctly", func(t *testing.T) {
		appErr := NewBadRequest("Invalid Payload", nil)
		rr := httptest.NewRecorder()
		bytes, err := json.Marshal(map[string]string{
			"error": "Invalid Payload",
		})
		assert.NoError(t, err)
		HandleError(rr, appErr)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		assert.Equal(t, string(bytes) + "\n", rr.Body.String())
	})

	t.Run("should handle regular error as internal server error", func(t *testing.T) {
		err := errors.New("Error")
		rr := httptest.NewRecorder()

		bytes, err := json.Marshal(map[string]string{
			"error": "Internal server error",
		})
		assert.NoError(t, err)

		HandleError(rr, err)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, string(bytes) + "\n", rr.Body.String())
	})
}
