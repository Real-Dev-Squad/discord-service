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
		assert.Equal(t, w.Code, 400)
	})
}

func TestNewUnauthorisedError(t *testing.T) {
	t.Run("should write unauthorised error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		Errors.NewUnauthorisedError(w, "test")
		assert.Equal(t, w.Code, 401)
	})
	t.Run("should write unauthorised error response with default message", func(t *testing.T) {
		w := httptest.NewRecorder()
		Errors.NewUnauthorisedError(w)
		assert.Equal(t, w.Code, 401)
		assert.Equal(t, w.Body.String(), fmt.Sprintln(`{"success": false, "message": "Unauthorized Access", "status": 401}`))
	})
}