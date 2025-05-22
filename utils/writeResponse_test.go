package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {

	t.Run("Should return error when json encoding fails", func(t *testing.T) {
		response := httptest.NewRecorder()
		err := WriteResponse(make(chan int), response)
		assert.Error(t, err)
	})

	t.Run("Should encode data in response body", func(t *testing.T) {
		response := httptest.NewRecorder()
		err := WriteResponse(map[string]interface{}{"message": "Hello, World!"}, response)
		assert.NoError(t, err)
		assert.Equal(t, `{"message":"Hello, World!"}` + "\n", response.Body.String())
	})

}
