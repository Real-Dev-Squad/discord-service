package utils

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDiscordResponse(t *testing.T) {
	t.Run("should write error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		Success.NewDiscordResponse(w, "test", make(chan int))
		assert.Equal(t, fmt.Sprintln(`{"success": false, "message": "Internal Server Error", "status": 500}`), w.Body.String())
	})

	t.Run("should write success response when data is nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		Success.NewDiscordResponse(w, "test", nil)
		assert.Equal(t, `{"success": true, "status": 200, "message": "test"}`, w.Body.String())
	})

	t.Run("should write success response when data is not nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		Success.NewDiscordResponse(w, "test", map[string]string{"value": "hello"})
		assert.Equal(t, `{"value":"hello"}`+"\n", w.Body.String())
	})
}
