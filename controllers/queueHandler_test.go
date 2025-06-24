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
	router:= httprouter.New()
	router.POST("/queue", controllers.QueueHandler)

	t.Run("should return internal server error when fails to read body", func(t *testing.T){
		req := httptest.NewRequest("POST", "/queue", &errorReader{})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("should return 200 when handler returns nil", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/queue", bytes.NewBuffer(make([]byte, 0)))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should fail if ExponentialBackoffRetry fails for listening command", func(t *testing.T) {
		config.AppConfig.MAX_RETRIES = 1
		originalFunc := utils.ExponentialBackoffRetry
		utils.ExponentialBackoffRetry = func(maxRetries int, operation func() error) error {
			return errors.New("error")
		}
		defer func() { utils.ExponentialBackoffRetry = originalFunc }()
		body := []byte(`{"CommandName": "listening", "MetaData": {"value": "true", "nickname" : "joy-gupta-1"}}`)
		req := httptest.NewRequest("POST", "/queue", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
