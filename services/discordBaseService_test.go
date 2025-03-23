package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/fixtures"
	service "github.com/Real-Dev-Squad/discord-service/services"
	_ "github.com/Real-Dev-Squad/discord-service/tests/setup"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestDiscordBaseService(t *testing.T) {
	t.Run("Should return 400 when request body is empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/", http.NoBody)
		service.DiscordBaseService(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 400 when request body is malformed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/", bytes.NewBuffer([]byte("malformed request")))
		service.DiscordBaseService(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 200 when request body is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		message := dtos.DiscordMessage{Type: discordgo.InteractionPing}
		jsonBytes, err := json.Marshal(message)
		assert.Nil(t, err)
		r, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonBytes))
		service.DiscordBaseService(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)

		assert.Nil(t, err)
		assert.Equal(t, float64(discordgo.InteractionResponsePong), response["type"])
	})
	t.Run("Should return 200 when request body is valid for interaction command", func(t *testing.T) {
		w := httptest.NewRecorder()
		jsonBytes, err := json.Marshal(fixtures.HelloCommand)
		assert.Nil(t, err)
		r, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonBytes))
		service.DiscordBaseService(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)

		assert.Nil(t, err)
		assert.Equal(t, float64(discordgo.InteractionResponseChannelMessageWithSource), response["type"])
	})

	t.Run("Should return 404 when interaction type is unknown", func(t *testing.T) {
		w := httptest.NewRecorder()
		message := dtos.DiscordMessage{}
		jsonBytes, err := json.Marshal(message)
		assert.Nil(t, err)
		r, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonBytes))
		service.DiscordBaseService(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
