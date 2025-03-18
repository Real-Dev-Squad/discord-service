package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
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

type TestResponse struct {
	Success bool `json:"success"`
	Data    []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
}

func TestMakeAPICall(t *testing.T) {
	t.Run("should make API call successfully", func(t *testing.T) {
		mockServer := MakeMockServer(`{"success": true, "data" : [{"name" : "joy", "age" : 20}]}`, http.StatusOK)
		defer mockServer.Close()

		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    http.MethodPost,
			URL:       mockServer.URL,
		}
		body := map[string]string{"key": "value"}
		result := &TestResponse{}
		err := wb.MakeAPICall(body, result)

		assert.NoError(t, err)
		assert.Equal(t, true, result.Success)
		assert.Equal(t, "joy", result.Data[0].Name)
		assert.Equal(t, 20, result.Data[0].Age)
	})

	t.Run("should make API call successfully", func(t *testing.T) {
		mockServer := MakeMockServer(`{"success": true, "data" : [{"name" : "joy", "age" : 20}]}`, http.StatusOK)
		defer mockServer.Close()

		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    http.MethodPost,
			URL:       mockServer.URL,
		}
		body := map[string]string{"key": "value"}
		result := &TestResponse{}
		err := wb.MakeAPICall(body, result)

		assert.NoError(t, err)
		assert.Equal(t, true, result.Success)
		assert.Equal(t, "joy", result.Data[0].Name)
		assert.Equal(t, 20, result.Data[0].Age)
	})

	t.Run("should return error if json.Marshal fails", func(t *testing.T) {
		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    http.MethodPost,
			URL:       "http://example.com/api",
		}
		body := make(chan int) // This will cause json.Marshal to fail
		err := wb.MakeAPICall(body, nil)

		assert.Error(t, err)
	})

	t.Run("should return error if http.NewRequest fails", func(t *testing.T) {
		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    "invalid method", // This will cause http.NewRequest to fail
			URL:       "http://example.com/api",
		}
		body := map[string]string{"key": "value"}
		err := wb.MakeAPICall(body, nil)

		assert.Error(t, err)
	})

	t.Run("should handle error from client.Do", func(t *testing.T) {

		mockServer := MakeMockServer(`{"success": true, "data" : [{"name" : "joy", "age" : 20}]}`, http.StatusOK)
		defer mockServer.Close()
		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    http.MethodPost,
			URL:       mockServer.URL,
		}
		originalTimeout := config.AppConfig.TIMEOUT
		config.AppConfig.TIMEOUT = 2
		defer func() { config.AppConfig.TIMEOUT = originalTimeout }()
		body := map[string]string{"key": "value"}
		result := &TestResponse{}
		err := wb.MakeAPICall(body, result)
		assert.Error(t, err)
	})

	t.Run("should handle server error responses", func(t *testing.T) {
		result := &TestResponse{}
		mockServer := MakeMockServer(`{"success":false}`, http.StatusInternalServerError)
		defer mockServer.Close()

		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    http.MethodPost,
			URL:       mockServer.URL,
		}
		body := map[string]string{"key": "value"}
		err := wb.MakeAPICall(body, result)

		assert.NoError(t, err)
		assert.Equal(t, false, result.Success)
	})

	t.Run("should handle server error responses", func(t *testing.T) {
		result := &TestResponse{}
		mockServer := MakeMockServer(`{"success":false}`, http.StatusInternalServerError)
		defer mockServer.Close()

		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    http.MethodPost,
			URL:       mockServer.URL,
		}
		body := map[string]string{"key": "value"}
		err := wb.MakeAPICall(body, result)

		assert.NoError(t, err)
		assert.Equal(t, false, result.Success)
	})

	t.Run("should return error if response does not matches with the provided dto", func(t *testing.T) {
		result := &TestResponse{}
		mockServer := MakeMockServer(`{success:false}`, http.StatusInternalServerError)
		defer mockServer.Close()

		wb := &WebsiteBackend{
			AuthToken: nil,
			Method:    http.MethodPost,
			URL:       mockServer.URL,
		}
		body := map[string]string{"key": "value"}
		err := wb.MakeAPICall(body, result)
		assert.Error(t, err)
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

func MakeMockServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
}
