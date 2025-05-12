package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloResponse(t *testing.T) {
	t.Run("should return hello response", func(t *testing.T) {
		response := ResponseGenerator.HelloResponse("123")
		assert.Equal(t, response, "Hey there <@123>! Congratulations, you just executed your first slash command")
	})
}
