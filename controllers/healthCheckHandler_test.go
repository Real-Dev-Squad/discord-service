package controllers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Real-Dev-Squad/discord-service/controllers"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Run("should return 200 OK", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.HealthCheckHandler(w, r, nil)
		})

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response map[string]string
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])

		// Validate the timestamp format
		timestamp, exists := response["timestamp"]
		assert.True(t, exists, "timestamp field should exist")
		_, err = time.Parse(time.RFC3339, timestamp)
		assert.NoError(t, err, "timestamp should be in RFC3339 format")
	})

	t.Run("should return 500 Internal Server Error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.HealthCheckHandler(w, r, nil)
		})

		originalFunc := utils.Encode
		utils.Encode = func(w io.Writer, data interface{}) error {
			return errors.New("simulated error")
		}
		defer func() {
			utils.Encode = originalFunc
		}()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

	})
}
