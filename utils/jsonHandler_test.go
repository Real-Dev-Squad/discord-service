package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonHandler(t *testing.T) {
	t.Run("should return json response", func(t *testing.T) {
		res, err:= Json.ToJson(map[string]string{"name": "John", "age": "20"})
		assert.Equal(t, nil, err)
		assert.Equal(t, `{"age":"20","name":"John"}`, res)
	})

	t.Run("should return error when fails to marshal data", func(t *testing.T) {
		res, err:= Json.ToJson(make(chan int))
		assert.Equal(t, "", res)
		assert.Equal(t, "json: unsupported type: chan int", err.Error())
	})
}
