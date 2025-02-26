package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/fixtures"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/Real-Dev-Squad/discord-service/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestHelloService(t *testing.T) {

	t.Run("should return a success response with a message", func(t *testing.T) {
		w := httptest.NewRecorder()
		jsonBytes, _ := json.Marshal(fixtures.HelloCommand)
		r, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonBytes))

		CS := CommandService{
			discordMessage: fixtures.HelloCommand,
		}
		CS.Hello(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response discordgo.InteractionResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, discordgo.InteractionResponseChannelMessageWithSource, response.Type)
		assert.Equal(t, utils.ResponseGenerator.HelloResponse(fixtures.HelloCommand.Member.User.ID), response.Data.Content)
	})
}
