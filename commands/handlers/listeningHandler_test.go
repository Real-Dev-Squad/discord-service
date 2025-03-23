package handlers

import (
	"errors"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/models"
	_ "github.com/Real-Dev-Squad/discord-service/tests/setup"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/stretchr/testify/assert"
)

type MockCommandHandler struct {
	discordMessage *dtos.DataPacket
}

func TestListeningHandler(t *testing.T) {

	t.Run("should update nickname with prefix and suffix if value is true", func(t *testing.T) {

		dataPacket := &dtos.DataPacket{
			UserID: "userID",
			MetaData: map[string]string{
				"nickname": "testNick",
				"value":    "true",
			},
		}

		handler := &CommandHandler{discordMessage: dataPacket}
		err := handler.listeningHandler()
		assert.Error(t, err)
	})

	t.Run("should update nickname without prefix and suffix if value is false", func(t *testing.T) {

		dataPacket := &dtos.DataPacket{
			UserID: "userID",
			MetaData: map[string]string{
				"nickname": utils.NICKNAME_PREFIX + "testNick" + utils.NICKNAME_SUFFIX,
				"value":    "false",
			},
		}

		handler := &CommandHandler{discordMessage: dataPacket}
		err := handler.listeningHandler()
		assert.Error(t, err)
	})

	t.Run("should return error if CreateSession fails", func(t *testing.T) {
		dataPacket := &dtos.DataPacket{
			UserID: "userID",
			MetaData: map[string]string{
				"nickname": "testNick",
				"value":    "true",
			},
		}

		originalFunc := models.CreateSession
		defer func() { models.CreateSession = originalFunc }()
		models.CreateSession = func() (*models.SessionWrapper, error) {
			return nil, errors.New("Failed")
		}
		handler := &CommandHandler{discordMessage: dataPacket}
		err := handler.listeningHandler()
		assert.Error(t, err)
		assert.Equal(t, "Failed", err.Error())
	})
	t.Run("should return error if UpdateNickName fails", func(t *testing.T) {

		dataPacket := &dtos.DataPacket{
			UserID: "userID",
			MetaData: map[string]string{
				"nickname": "testNick",
				"value":    "true",
			},
		}

		handler := &CommandHandler{discordMessage: dataPacket}
		err := handler.listeningHandler()
		assert.Error(t, err)
		assert.Equal(t, "HTTP 401 Unauthorized, {\"message\": \"401: Unauthorized\", \"code\": 0}", err.Error())
	})
}
