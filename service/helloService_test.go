package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/fixtures"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestHelloService(t *testing.T) {
	t.Run("should return a success response with a message", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		CS.discordMessage = fixtures.HelloCommand
		CS.HelloService(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		var response discordgo.InteractionResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, discordgo.InteractionResponseChannelMessageWithSource, response.Type)
		assert.Equal(t, "Hey there! Congratulations, you just executed your first slash command", response.Data.Content)
	})
}
