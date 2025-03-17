package config

import (
	"testing"

	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestValidateEnvironment(t *testing.T) {
	t.Run("should return development when environment is invalid", func(t *testing.T) {
		invalidEnv := Environment("invalid")
		validatedEnv := invalidEnv.Validate()
		assert.Equal(t, Development, validatedEnv)

	})

	t.Run("should return the same environment when it is valid", func(t *testing.T) {
		validEnv := Production
		validatedEnv := validEnv.Validate()
		assert.Equal(t, Production, validatedEnv)
	})

}
