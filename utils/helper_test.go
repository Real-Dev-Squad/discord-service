package utils

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

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

		wb := &HTTPClientService{
			AuthToken: "",
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
		wb := &HTTPClientService{
			AuthToken: "",
			Method:    http.MethodPost,
			URL:       "http://example.com/api",
		}
		body := make(chan int) // This will cause json.Marshal to fail
		err := wb.MakeAPICall(body, nil)

		assert.Error(t, err)
	})

	t.Run("should return error if http.NewRequest fails", func(t *testing.T) {
		wb := &HTTPClientService{
			AuthToken: "",
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
		wb := &HTTPClientService{
			AuthToken: "",
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

		wb := &HTTPClientService{
			AuthToken: "",
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

		wb := &HTTPClientService{
			AuthToken: "",
			Method:    http.MethodPost,
			URL:       mockServer.URL,
		}
		body := map[string]string{"key": "value"}
		err := wb.MakeAPICall(body, result)
		assert.Error(t, err)
	})
}

func TestPrepareRequest(t *testing.T) {

	t.Run("should prepare request with valid body and headers", func(t *testing.T) {
		hcp := &HTTPClientService{
			AuthToken: "test-token",
			Method:    http.MethodPost,
			URL:       "http://example.com",
		}

		body := map[string]string{"key": "value"}
		req, err := hcp.prepareRequest(body)

		assert.NoError(t, err)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", req.Header.Get("Authorization"))
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "http://example.com", req.URL.String())
	})

	t.Run("should return error if json.Marshal fails", func(t *testing.T) {
		hcp := &HTTPClientService{
			Method: http.MethodPost,
			URL:    "http://example.com",
		}

		body := make(chan int) // This will cause json.Marshal to fail
		_, err := hcp.prepareRequest(body)

		assert.Error(t, err)
	})

	t.Run("should return error if http.NewRequest fails", func(t *testing.T) {
		mockServer := MakeMockServer(`{"success": true, "data" : [{"name" : "joy", "age" : 20}]}`, http.StatusOK)
		defer mockServer.Close()

		hcp := &HTTPClientService{
			Method: "invalid-method",
			URL:    mockServer.URL,
		}

		body := interface{}(nil)
		err := hcp.MakeAPICall(body, nil)

		assert.Error(t, err)
	})
}

func TestReadResponseBody(t *testing.T) {
	hcp := &HTTPClientService{}
	respBody := `{"key": "value"}`
	resp := &http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte(respBody))),
	}

	var result map[string]string
	err := hcp.readResponseBody(resp, &result)

	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func MakeMockServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}))
}
