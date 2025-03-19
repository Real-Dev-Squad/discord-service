package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func TestNewBadRequestError(t *testing.T) {
	response := httptest.NewRecorder()
	message := "Bad Request"
	Errors.NewBadRequestError(response, message)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	r := &result{}
	err := json.Unmarshal(response.Body.Bytes(), r)
	assert.NoError(t, err)
	assert.Equal(t, message, r.Message)
}

func TestNewUnauthorisedError(t *testing.T) {
	t.Run("should return 401 Unauthorised with default message", func(t *testing.T) {
		response := httptest.NewRecorder()
		Errors.NewUnauthorisedError(response)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		message := "Unauthorized Access"
		r := &result{}
		err := json.Unmarshal(response.Body.Bytes(), r)
		assert.NoError(t, err)
		assert.Equal(t, message, r.Message)
	})
	t.Run("should return 401 Unauthorised with provided error message", func(t *testing.T) {
		response := httptest.NewRecorder()
		Errors.NewUnauthorisedError(response, "test error")
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		r := &result{}
		err := json.Unmarshal(response.Body.Bytes(), r)
		assert.NoError(t, err)
		assert.Equal(t, "test error", r.Message)
	})
}
func TestNewInternalError(t *testing.T) {
	t.Run("should return 500 Internal Server Error with default message", func(t *testing.T) {
		response := httptest.NewRecorder()
		Errors.NewInternalError(response)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		message := "Internal Server Error"
		r := &result{}
		err := json.Unmarshal(response.Body.Bytes(), r)
		assert.NoError(t, err)
		assert.Equal(t, message, r.Message)
	})
	t.Run("should return 500 Internal Server Error with provided error message", func(t *testing.T) {
		response := httptest.NewRecorder()
		Errors.NewInternalError(response, "test error")
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		r := &result{}
		err := json.Unmarshal(response.Body.Bytes(), r)
		assert.NoError(t, err)
		assert.Equal(t, "test error", r.Message)
	})
}
