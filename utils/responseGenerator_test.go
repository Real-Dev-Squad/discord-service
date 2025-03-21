package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloResponse(t *testing.T) {
	t.Run("should return hello response", func(t *testing.T) {
		userId := "123"
		response := ResponseGenerator.HelloResponse(userId)
		expectedResponse := fmt.Sprintf("Hey there <@%s>! Congratulations, you just executed your first slash command", userId)
		assert.Equal(t, expectedResponse, response)
	})
}
