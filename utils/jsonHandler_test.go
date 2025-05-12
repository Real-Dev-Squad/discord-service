package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonHandler(t *testing.T) {
	t.Run("should return json response", func(t *testing.T) {
		assert.Equal(t, `{"age":"20","name":"John"}`, Json.ToJson(map[string]string{"name": "John", "age": "20"}))
	})

	t.Run("should return empty string when fails to marshal data", func(t *testing.T) {
		assert.Equal(t, "", Json.ToJson(make(chan int)))
	})
}
