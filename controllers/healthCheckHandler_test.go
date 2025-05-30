package controllers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Real-Dev-Squad/discord-service/controllers"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.HealthCheckHandler(w, r, nil)
	})
	
	t.Run("should return 200 status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])

		// Validate the timestamp format
		timestamp, exists := response["timestamp"]
		assert.True(t, exists, "timestamp field should exist")
		_, err = time.Parse(time.RFC3339, timestamp)
		assert.NoError(t, err, "timestamp should be in RFC3339 format")
	})

	t.Run("should return error response if json encoding fails", func(t *testing.T) {
		w := httptest.NewRecorder()

		originalFunc := utils.WriteResponse
		utils.WriteResponse = func(data interface{}, response http.ResponseWriter) error {
			return errors.New("test error")
		}
		defer func() { utils.WriteResponse = originalFunc }()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		res, _ := utils.Json.ToJson(utils.ErrorResponse{
			Success: false,
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		assert.Equal(t, fmt.Sprintln(res), w.Body.String())
	})
}
