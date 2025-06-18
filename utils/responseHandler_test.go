package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSONResponse(t *testing.T) {
	t.Run("should return success response with 200 status code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		data := map[string]string{"status": "ok"}
		status := http.StatusOK
		bytes, err:= json.Marshal(data)
		assert.NoError(t, err)
		WriteJSONResponse(rr, status, data)
		assert.Equal(t, status, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, string(bytes), rr.Body.String())
	})

	t.Run("should return response with 400 status code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		data := map[string]string{"message": "Invalid data"}
		status := http.StatusBadRequest
		WriteJSONResponse(rr, status, data)
		bytes, err := json.Marshal(data)
		assert.NoError(t, err)
		assert.Equal(t, status, rr.Code)
		assert.Equal(t, string(bytes), rr.Body.String())
	})

	t.Run("should have empty body when fails to marshal data", func(t *testing.T){
		rr := httptest.NewRecorder();
		bytes, err:= json.Marshal(map[string]string{
			"error": "Internal Server Error",
		})
		assert.NoError(t, err)
		WriteJSONResponse(rr, http.StatusAccepted, make(chan int))
		assert.Equal(t, string(bytes) + "\n", rr.Body.String())
	})
}