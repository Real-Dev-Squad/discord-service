package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonHandler(t *testing.T) {
	t.Run("should return json response", func(t *testing.T) {
		response := Json.ToJson(map[string]string{"name": "John", "age": "20"})
		assert.Equal(t, response, `{"age":"20","name":"John"}`)
	})

	t.Run("should return empty string when fails to marshal data", func(t *testing.T) {
		response := Json.ToJson(make(chan int))
		assert.Equal(t, response, "")
	})
}
