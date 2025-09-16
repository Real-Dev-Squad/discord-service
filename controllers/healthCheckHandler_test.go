package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Real-Dev-Squad/discord-service/controllers"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Run("should return 200 status code", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.HealthCheckHandler(w, r, nil)
		})

		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotNil(t, w.Body.String())
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		res := map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "ok", res["status"])

		timestamp, exists := res["timestamp"]
		assert.True(t, exists, "timestamp field should exist")
		_, err = time.Parse(time.RFC3339, timestamp)
		assert.NoError(t, err, "timestamp should be in RFC3339 format")
	})
}
