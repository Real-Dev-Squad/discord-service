package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoffRetry_Success(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := ExponentialBackoffRetry(5, operation)
	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
}

func TestExponentialBackoffRetry_Failure(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return errors.New("permanent error")
	}

	err := ExponentialBackoffRetry(3, operation)
	assert.Error(t, err)
	assert.Equal(t, 3, attempts)
}

func TestExponentialBackoffRetry_NoRetries(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return errors.New("error")
	}

	err := ExponentialBackoffRetry(0, operation)
	assert.Nil(t, err)
	assert.Equal(t, 0, attempts)
}

func TestExponentialBackoffRetry_ImmediateSuccess(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return nil
	}
	err := ExponentialBackoffRetry(5, operation)
	assert.NoError(t, err)
	assert.Equal(t, 1, attempts)
}

func TestMakeAPICall(t *testing.T) {
	t.Run("should make API call successfully", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"success": true}`)
		}))
		defer mockServer.Close()

		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    http.MethodPost,
			URL:       mockServer.URL,
		}
		body := map[string]string{"key": "value"}
		resp, err := wb.MakeAPICall(body)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("should return error if json.Marshal fails", func(t *testing.T) {
		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    http.MethodPost,
			URL:       "http://example.com/api",
		}
		body := make(chan int) // This will cause json.Marshal to fail
		resp, err := wb.MakeAPICall(body)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error if http.NewRequest fails", func(t *testing.T) {
		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    "invalid method", // This will cause http.NewRequest to fail
			URL:       "http://example.com/api",
		}
		body := map[string]string{"key": "value"}
		resp, err := wb.MakeAPICall(body)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestPrepareHeaders(t *testing.T) {
	t.Run("should set Content-Type header", func(t *testing.T) {
		wb := &WebsiteBackend{
			AuthToken: nil,
		}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		assert.NoError(t, err)
		wb.PrepareHeaders(req)

		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	})

	t.Run("should set Authorization header if AuthToken is provided", func(t *testing.T) {
		token := "test-token"
		wb := &WebsiteBackend{
			AuthToken: &token,
		}
		req, err := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		assert.NoError(t, err)
		wb.PrepareHeaders(req)

		assert.Equal(t, "Bearer test-token", req.Header.Get("Authorization"))
	})
}
