package utils

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/stretchr/testify/assert"
)

func TestNewDiscordResponse(t *testing.T) {
	t.Run("should return error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		Success.NewDiscordResponse(w, "test", make(chan int))
		res, _:= Json.ToJson(dtos.Response{
			Success: false,
			Message: "Internal Server Error",
			Status:  500,
		})
		assert.Equal(t, fmt.Sprintln(res), w.Body.String())
	})

	t.Run("should return success response when data is nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		Success.NewDiscordResponse(w, "test", nil)
		res, _:= Json.ToJson(dtos.Response{
			Success: true,
			Message: "test",
			Status:  200,
		})
		assert.Equal(t, res, w.Body.String())
	})

	t.Run("should return success response when data is not nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		Success.NewDiscordResponse(w, "test", map[string]string{"value": "hello"})
		assert.Equal(t, `{"value":"hello"}`+"\n", w.Body.String())
	})
}
