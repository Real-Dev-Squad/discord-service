package commandHandler

import (
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	t.Run("should return listeningHandler for 'listening' command", func(t *testing.T) {
		dataPacket := &dtos.DataPacket{
			CommandName: "listening",
		}
		data, err := dataPacket.ToByte()
		assert.NoError(t, err)

		handler := MainHandler(data)
		assert.NotNil(t, handler)
	})
	t.Run("should return nil for invalid data", func(t *testing.T) {
		invalidData := []byte("invalid data")

		handler := MainHandler(invalidData)
		assert.Nil(t, handler)
	})
}
