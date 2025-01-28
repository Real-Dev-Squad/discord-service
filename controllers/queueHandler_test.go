package controllers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/controllers"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}
func TestQueueHandler(t *testing.T) {

	router := httprouter.New()
	router.POST("/queue", controllers.QueueHandler)
	t.Run("should return 200 OK and log the request body", func(t *testing.T) {
		body := []byte(`{"message": "test message"}`)
		req, err := http.NewRequest("POST", "/queue", bytes.NewBuffer(body))
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	})
	t.Run("should be able to execute listening command", func(t *testing.T) {
		body := []byte(`{"CommandName": "listening"}`)
		req, err := http.NewRequest("POST", "/queue", bytes.NewBuffer(body))
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	})

	t.Run("should fail if ExponentialBackoffRetry fails for listening command", func(t *testing.T) {
		config.AppConfig.MAX_RETRIES = 1
		originalFunc := utils.ExponentialBackoffRetry
		utils.ExponentialBackoffRetry = func(maxRetries int, operation func() error) error {
			return errors.New("error")
		}
		defer func() { utils.ExponentialBackoffRetry = originalFunc }()
		body := []byte(`{"CommandName": "listening"}`)
		_, err := http.NewRequest("POST", "/queue", bytes.NewBuffer(body))
		assert.NoError(t, err)
	})

	t.Run("should return 500 Internal Server Error if payload is unable to be decoded", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/queue", &errorReader{})
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
