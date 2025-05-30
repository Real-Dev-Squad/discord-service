package utils

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBadRequestError(t *testing.T) {
	t.Run("should return bad request error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		Errors.NewBadRequestError(w, "test")
		assert.Equal(t, 400, w.Code)
		error, _:= Json.ToJson(ErrorResponse{
			Success: false,
			Message: "test",
			Status:  400,
		})
		assert.Equal(t, fmt.Sprintln(error), w.Body.String())
	})
}

func TestNewUnauthorisedError(t *testing.T) {
	t.Run("should return unauthorised error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		Errors.NewUnauthorisedError(w, "test")
		assert.Equal(t, 401, w.Code)
		error, _:= Json.ToJson(ErrorResponse{
			Success: false,
			Message: "test",
			Status:  401,
		})
		assert.Equal(t, fmt.Sprintln(error), w.Body.String())
	})
	t.Run("should return unauthorised error response with default message", func(t *testing.T) {
		w := httptest.NewRecorder()
		Errors.NewUnauthorisedError(w)
		assert.Equal(t, 401, w.Code)
		error, _:= Json.ToJson(ErrorResponse{
			Success: false,
			Message: "Unauthorized Access",
			Status:  401,
		})
		assert.Equal(t, fmt.Sprintln(error), w.Body.String())
	})
}