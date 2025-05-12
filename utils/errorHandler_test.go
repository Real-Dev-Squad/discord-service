package utils

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBadRequestError(t *testing.T) {
	t.Run("should write bad request error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		Errors.NewBadRequestError(w, "test")
		assert.Equal(t, 400, w.Code)
	})
}

func TestNewUnauthorisedError(t *testing.T) {
	t.Run("should write unauthorised error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		Errors.NewUnauthorisedError(w, "test")
		assert.Equal(t, 401, w.Code)
	})
	t.Run("should write unauthorised error response with default message", func(t *testing.T) {
		w := httptest.NewRecorder()
		Errors.NewUnauthorisedError(w)
		assert.Equal(t, 401, w.Code)
		assert.Equal(t, fmt.Sprintln(`{"success": false, "message": "Unauthorized Access", "status": 401}`), w.Body.String())
	})
}